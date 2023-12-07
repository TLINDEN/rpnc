-- simple function, return the lower number of the two operands
function lower(a,b)
    if a < b then
        return a
    else
        return b
    end
end

function init()
    -- expects 2 args
    register("lower", 2, "lower")
end  
