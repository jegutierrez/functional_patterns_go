# fuctntional_patterns_go
Recopilacion de algunos patrones que he identificado en aplicaciones en las que he trabajado.

### Entonces, ¿programación funcional en Go?

Go nos provee la capacidad de usar funciones como ciudadanos de primera clase, es decir, podemos utilizar las funciones para pasarlas como parámetro o retornarlas como valor de otra función (esto es conocido como funciones de orden superior), sin embargo, como todo en la ingenieria de software, es acerca de tradeoffs, lo recomendado es siempre preferir la claridad del código antes que cualquier patron o paradicma, por muy cool que sea.

En el siguiente post vamos a describir una serie de técnicas de programación funcional, aprovechando esta capacidad de trabajar con funciones como valor.

### Recursividad & benchmarks

Cuando hablamos de programación funcional tenemos dos requisitos importantes, trabajar con funciones y no mutar estado. Y una de las primeras tecnicas que nos viene a la cabeza es la recursividad, que es la capacidad de una función de llamarse a si misma. Esto es necesario porque es la manera de hacer loops en FP.

Por lo tanto lo primero que no podemos usar si queremos mantener inmutabilidad en nuestra aplicacion es la instrucción `for`, ya que dentro vamos mutando una variable que toma diferente valor en cada iteración.

Otra de las cosas que acostumbramos a usar en Go son estructuras de datos a las que vamos appendeando elementos, esto tampoco podríamos, aunque en el caso del `slice` usando la función `append()` no estamos mutando ya que nos retorna un nuevo slice cada vez que se llama.

¿Es buena idea esto en Go?

---
Mostrar código de recursividad vs for y estructuras mutables y benchmarks.
---

Despues de haber visto los ejemplos nos podemos dar cuenta que no es tan buena idea usar recursividad sobre `for` en Go, ya que el `for` está muy optimizado, no tenemos problema de que crezca mucho el stack de memoria, ademas de ser la manera mas idiomatica de resolver casi todos los problemas en el ecosistema, sin embargo, como dije antes, lo que tenemos que priorizar siempre es la claridad del código y si para el equipo una solución funcional resulta mas clara, pues entonces es la manera correcta.

Ademas cabe destacar que el problema de performance que vimos en los benchmarks, va a ser insignificante en el 99% de los casos, normalmente vamos a tener cullos de botella en otro lugar fuera del código, sobre todo si tenemos llamadas a servicios remotos, bases de datos, etc.

### Higher order functions

Mostrar como podemos hacer una funcion que retorna un http handler y poder hacer cosas antes de retornar el handler

### Closures

Pasando closures a funciones en los tests podemos mockear funciones que nos interesa testear, ademas de poder hacer asserts dentro del closure conservando el scope de `t *testing.T`

Usando la tecnica de closures podemos crear una funcion para filtrar por monto en dos tipos diferentes, la idea es declarar un slice fuera del closure, pasar la referencia a la funcion filter y luego appendear dentro al slice.

---
ejemplo de una funcion filter por monto en dos tipos diferentes
---

Explicar que de alguna manera esto suple el uso de generics en Go para este caso en particular, si queda tiempo dar un pequeño overview de porque es dificil implementar generics en Go

Explicar que esto es posible debido a que tenemos un dato en comun, en caso de los slices es el indice tipo `int`, pero cabe destacar que esto mismo no sería posible usando maps, debido a que las keys pueden ser de cualquier tipo.

---
ejemplo de que no funciona con map
---

### Partial Aplication

Para los test utilizamos esta tecnica, de manera de aplicar el context de 

### Concurrencia y funciones

Aqui quiero destacar y mostrar el ejemplo de hacer llamadas a multiples APIs como un caso de la vida real.

Mostrar el contraste con promesas en javascript y Futures en java, en Go no tenemos algo igual en la standard library, podemos usar una libreria para manejar promesas o futuros, pero en mi opinion, eso no sería el estilo de Go.

Dentro de Go tenemos las Goruoutines y los channels para manejar muy elegantemente la mayoria de los casos de concurrencia que se nos puedan presentar.

Mostrar ejemplos de código y benchmarks:

----
Hacer llamadas a api bloqueantes
----

----
Hacer llamadas a api usando wait groups y closures
----

----
Hacer llamadas a api usando channels y closures
----

Mostrar como aprovechando la capacidad de pasar funciones podemos hacer patrones de concurrencia muy elegantes, ademas de poder hacer genericas algunas funciones dentro de nuestro código

### Conclusiones