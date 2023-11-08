-- simple function, return the lower number of the two operands
function lower(a,b)
    if a < b then
        return a
    else
        return b
    end
end

-- calculate parallel resistance. Batch  function (registered with -1,
-- see below). Takes a table as parameter.
--
-- Formula: 1/( (1/R1) + (1/R2) + ...)
function parallelresistance(list)
    sumres = 0
    
    for i, value in ipairs(list) do
        sumres = sumres + 1 / value
    end

    return 1 / sumres
end

-- converter example
function inch2centimeter(inches)
    return inches * 2.54
end

function init()
    -- expects 2 args
    register("lower", 2, "lower")

    -- expects a list of all numbers on the stack, batch mode
    register("parallelresistance", -1, "parallel resistance")

    -- expects 1 arg, but doesn't pop()
    register("inch2centimeter", 0)
end  
