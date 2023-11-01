function add(a,b)
  return a + b
end

function test(a)
  return a
end  

function parallelresistance(a,b)
  return 1.0 / (a * b)
end

function init()
  register("add", 2, "addition")
  register("test", 1, "test")
  register("parallelresistance", 2, "parallel resistance")
end  
