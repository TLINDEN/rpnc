function add(a,b)
  return a + b
end

function parallelresistance(a,b)
  return 1.0 / (a * b)
end

function init()
  RegisterFuncTwoArg("add")
  RegisterFuncTwoArg("parallelresistance")
end  
