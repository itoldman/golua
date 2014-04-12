#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "lua.h"
#include "lauxlib.h"

#ifdef NDEBUG
  #define DEBUG(x,args...) ((void) (0))
  #define DEBUG_MSG(msg,buf) ((void) (0))
#else
  #define DEBUG(x,args...) debug_error(x, ##args);
  #define DEBUG_MSG(msg,buf) ei_show_recmsg(stderr, &msg, (char *)&buf);
#endif

#define MAXATOMLEN 255
int fd = -1;

// Dump an error message to stderr.
void debug_error(const char *fmt, ...)
{
	va_list ap;
	va_start(ap, fmt);
	vfprintf(stderr, fmt, ap);
	fputs("\n", stderr);
	va_end(ap);
}

// Dump a fatal error to stderr and exit the program.
void fatal_error(const char *fmt, ...)
{
	va_list ap;
	va_start(ap, fmt);
	vfprintf(stderr, fmt, ap);
	fputs("\n", stderr);
	va_end(ap);
	exit(1);
}

// Dump a lua stack element, without recursion
void dump_element(lua_State* L, int idx, int depth)
{
	char d[depth+1];
	int i;
	for (i=0; i<depth; ++i) d[i] = ' ';
	d[depth] = '\0';

	switch(lua_type(L, idx))
	{
		case LUA_TNIL:
			DEBUG("%s[%d] (%s) %s",d,idx,"nil","nil");
			break;
		case LUA_TBOOLEAN:
			DEBUG("%s[%d] (%s) %s",d,idx,"boolean",lua_toboolean(L, idx) ? "true" : "false" );
			break;
		case LUA_TNUMBER:
			DEBUG("%s[%d] (%s) %f",d,idx,"number",lua_tonumber(L, idx));
			break;
		case LUA_TSTRING:
			DEBUG("%s[%d] (%s) %s",d,idx,"string",lua_tostring(L, idx));
			break;
		case LUA_TTABLE:
			DEBUG("%s[%d] (%s)",d,idx,"table");
			break;
		case LUA_TLIGHTUSERDATA:
			DEBUG("%s[%d] (%s)",d,idx,"light user data");
			break;
		case LUA_TFUNCTION:
			DEBUG("%s[%d] (%s)",d,idx,"function");
			break;
		case LUA_TUSERDATA:
			DEBUG("%s[%d] (%s)",d,idx,"user data");
			break;
		case LUA_TTHREAD:
			DEBUG("%s[%d] (%s)",d,idx,"thread");
			break;
		case LUA_TNONE:
			DEBUG("%s[%d] (%s)",d,idx,"none");
			break;
		default:
			break;
	}
}

// Dump a lua stack element, with recursion
void dump_element_recurse(lua_State* L, int idx, int depth)
{
	dump_element(L, idx, depth);
	if (lua_istable(L, idx)) {
		int n = lua_gettop(L);
		lua_pushvalue(L, idx); // push table on top
		dump_element_recurse(L, n+1, depth+1);
		lua_pop(L, 1); // pop table from tack
		assert( lua_gettop(L) == n );
	}
}

// Dump lua stack
void dump_stack(lua_State* L)
{
	DEBUG("=== STACK DUMP");
	int idx;
	for (idx=1; idx<=lua_gettop(L); ++idx)
		dump_element(L, idx, 0);
	DEBUG("=== END DUMP\n");
}

// Dump a lua variable, without recursion
void dump_var(lua_State* L, const char* var)
{
	DEBUG("== VAR DUMP");
	int n = lua_gettop(L);
	lua_getglobal(L, var);
	assert( lua_gettop(L) > 0 );
	dump_element(L, n+1, 0);
	lua_pop(L, 1);
	assert( lua_gettop(L) == n );
	DEBUG("== END DUMP\n");
}

// Dump a lua variable, with recursion
void dump_var_recurse(lua_State* L, const char* var)
{
	DEBUG("== VAR DUMP");
	int n = lua_gettop(L);
	lua_getglobal(L, var);
	dump_element_recurse(L, n+1, 0);
	lua_pop(L, 1);
	assert( lua_gettop(L) == n );
	DEBUG("== END DUMP\n");
}

/**
 * Dynamically cast a lua userdata type to the go struct of the given type name.
 *
 * @param L the lua state.
 * @param ud the lua stack index of the value to 'dynamic_cast'.
 * @param tname the expected type name.
 *
 * @return the userdata if type checking succeed, NULL pointer otherwise.
 */
void* luaL_castudata(lua_State* L, int ud, const char* tname)
{
	void *p = lua_touserdata(L, ud);
	if (p != NULL && lua_getmetatable(L, ud)) {		// does it have a metatable? 
		lua_getfield(L, LUA_REGISTRYINDEX, tname);  // get correct metatable 
		if (lua_rawequal(L, -1, -2)) {				// does it have the correct mt? 
			lua_pop(L, 2);							// remove both metatables 
			return p;
		}
	}
	lua_pop(L, 2); // remove both metatables
	return NULL;
}

