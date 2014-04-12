--[[
This is my first test lua program
Why to use it? I don't know yet.

--]]
width = 200 
height = 300
background = {r=0.30, g=0.10, b=0} 
array = {a, b, c}

a = {}
a.x = 10                    -- same as a["x"] = 10
print(a.x)                  -- same as print(a["x"])
print(a.y) 

function calculate(a, b)
	--(33 * a / b) ^ 15 * a + b
	return a + b
end
