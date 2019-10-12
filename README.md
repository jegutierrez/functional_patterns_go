# Functional patterns en Go
Recopilación de algunos patrones útiles que he identificado en el dia a dia trabajando con Go.

### Entonces, ¿programación funcional en Go?

Go nos provee la capacidad de usar funciones como ciudadanos de primera clase, es decir, podemos utilizar las funciones para pasarlas como parámetro o retornarlas como valor de otra función (esto es conocido como funciones de orden superior), sin embargo, como todo en software, es acerca de tradeoffs, lo recomendado es siempre preferir la claridad del código antes que cualquier patrón o paradigma.

En el siguiente post vamos a describir una serie de técnicas de programación funcional, aprovechando esta capacidad de trabajar con funciones como valor.

### Que implica hacer programación funcional

La programación funcional nos trae dos grandes restricciones:

* Usar funciones para todo: hacer todas las operaciones con funciones.
* No mutar estado: no mutar valores una vez declarados, no tener estructuras de dato mutables, no tener side effects dentro de nuestras funciones.

Algunos de los conceptos que nos sirve tener en cuenta para identificar los patrones funcionales son los siguientes:

### Recursividad & benchmarks

Cuando hablamos de programación funcional tenemos dos requisitos importantes, trabajar con funciones y no mutar estado. Y una de las primeras técnicas que nos viene a la cabeza es la recursividad, que es la capacidad de una función de llamarse a sí misma. Esto es necesario porque es la manera de hacer loops en FP.

Por lo tanto lo primero que no podemos usar si queremos mantener inmutabilidad en nuestra aplicación es *for loops*, ya que dentro vamos mutando una variable que toma diferente valor en cada iteración.

Otra de las cosas que acostumbramos a usar en Go son estructuras de datos a las que vamos agregando elementos, esto tampoco podríamos, aunque en el caso del *slice* usando la función `append()` no estamos mutando ya que nos retorna un nuevo slice cada vez que se llama.

¿Es buena idea no usar for loops en Go?

Veamos con un ejemplo y pruebas de performance:

Escribir una función que cuente cuántas maneras posibles de dar cambio hay para un monto dado con una lista de denominaciones de monedas.

Por ejemplo, hay 3 maneras de dar cambio para $4 con monedas de $1 y $2
* 1+1+1+1
* 1+1+2
* 2+2

Versión recursiva:

```go
func CoinsChangeRecursive(amount int, coins []int) int {
	if amount == 0 {
		return 1
	} else if amount > 0 && len(coins) > 0 {
		return CoinsChangeRecursive(amount-coins[0], coins) +
			CoinsChangeRecursive(amount, coins[1:])
	} else {
		return 0
	}
}
```

Versión usando for loop y una tabla para guardar los resultados de cada iteración:

```go
func CoinsChangeGoStyle(amount int, coins []int) int {
	var table = make([]int, amount+1, amount+1)
	table[0] = 1

	for i := 0; i < len(coins); i++ {
		for j := coins[i]; j <= amount; j++ {
			table[j] += table[j-coins[i]]
		}
	}
	return table[amount]
}
```

Benchmark

```go
func BenchmarkCoinsChangeRecursive(b *testing.B) {
	result := CoinsChangeRecursive(3000, []int{5, 10, 20, 50, 100, 200, 500})

	if result != 22481738 {
		b.Errorf("unspected result, want 22481738, got: %d", result)
	}
}

func BenchmarkCoinsChangeGoStyle(b *testing.B) {
	result := CoinsChangeGoStyle(3000, []int{5, 10, 20, 50, 100, 200, 500})

	if result != 22481738 {
		b.Errorf("unspected result, want 22481738, got: %d", result)
	}
}
```

Run the benchmark using *go test*
```sh
go test -bench=Coins
```

Results:
```sh
goos: darwin
goarch: amd64
BenchmarkCoinsChangeRecursive-4                1        21611593952 ns/op
BenchmarkCoinsChangeGoStyle-4           2000000000               0.00 ns/op
PASS
ok      _/Users/jegutierrez/Documents/projects/functional_patters_go/recursive  21.687s
```

Vemos que hay una diferencia considerable entre las dos versiones, la recursiva tarda un poco más de 21 segundos, esto es principalmente, debido a la cantidad de llamadas anidadas en el stack de ejecuciones (podríamos optimizar nuestra versión recursiva utilizando alguna técnica de memoization, pero solo busco mostrar que la recursividad puede tener un costo grande en algunos casos) y la versión que usa for loop tiene un tiempo cercano a cero.

Dejo todo el código y otro ejemplo utilizando tail recursividad, junto con los test y benchmarks aplicados [aqui](https://github.com/jegutierrez/functional_patterns_go/tree/master/recursion).

Después de haber visto los ejemplos y benchmarks nos podemos dar cuenta que no es tan buena idea usar recursividad en todos los casos en Go, ya que no tenemos problema de que crezca call el stack, además de ser la manera más idiomática de resolver casi todos los problemas en el ecosistema de Go, sin embargo, como dije antes, lo que tenemos que priorizar siempre es la claridad del código y si para el equipo una solución funcional resulta más clara, pues entonces es la manera correcta.

Como dato para pensar, Go ya es sumamente eficiente y en la mayoría de los casos el problema de performance que vimos en los benchmarks, va a ser insignificante, normalmente vamos a tener cuellos de botella en otro lugar fuera del código, sobre todo si tenemos llamadas a través de la red, manejo de archivos, bases de datos, etc.

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

Frecuentemente utilizamos técnicas funcionales para trabajar con concurrencia en Go como funciones lambda o closures. 

Vamos a mostrar el ejemplo muy común de hacer llamadas a múltiples APIs y luego unir los resultados, lo que buscamos es, hacer múltiples requests http en paralelo y luego esperar a que vuelvan las respuestas para procesarlas.

En Go no tenemos en la standard library algo parecido a las promesas en javascript o Futures en java, sin embargo utilizado goroutines, waitgroups y channel podemos construir nuestra propia manera de resolver promesas al estilo simple de Go.

Primero mostramos un ejemplo de 3 endpoint, con delays:

 /users 		-> 150 ms
 /balance 		-> 350 ms
 /user-debts 		-> 250 ms


1. Primero hacemos un cliente con los 3 request bloqueantes, sin paralelismo:

En este caso:
* En este caso hacemos 3 requests http y cada uno bloquea hasta completarse.
* Ignoramos los errores para hacer más concreto el ejemplo.
* Se unen las respuesta una vez que se completaron las 3.

```go
func GetUserStatusSync(serverURL, userID string) (UserStatus, error) {
	userResponse, _ := http.Get(fmt.Sprintf("%s/users/%s", serverURL, userID))
	balanceResponse, _ := http.Get(fmt.Sprintf("%s/balance/%s", serverURL, userID))
	debtsResponse, _ := http.Get(fmt.Sprintf("%s/user-debts/%s", serverURL, userID))

	var userInfo map[string]string
	unmarshalResponse(userResponse, &userInfo)
	var userBalance map[string]string
	unmarshalResponse(balanceResponse, &userBalance)
	var userDebts []map[string]string
	unmarshalResponse(debtsResponse, &userDebts)

	return UserStatus{
		ID:            userInfo["id"],
		Name:          userInfo["name"],
		BalanceAmount: userBalance["amount"],
		Debts:         userDebts,
	}, nil
}
```

2. Una llamada utilizando funciones anónimas, closures, goroutines y un waitgroup para esperar a las respuesta de los 3 endpoints.

En este caso:
* Declaramos 1 waitgroup para esperar a los 3 requests.
* Declaramos 3 variables `userResponse, balanceResponse, debtsResponse` de tipo `*http.Response` antes de las llamadas http.
* Cada request http lo hacemos dentro de una lambda y cada lambda se ejecuta en una goroutine diferente y marca cómo `Done()` en el waitgroup una vez que tiene la respuesta.
* Bloqueamos con el waitgroup hasta que se completen las 3 llamadas.
* Se unen las respuesta como en el caso anterior.

```go
func GetUserStatusAsyncWaitGroup(serverURL, userID string) (UserStatus, error) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(3)

	var userResponse, balanceResponse, debtsResponse *http.Response
	go func() {
		userResponse, _ = http.Get(fmt.Sprintf("%s/users/%s", serverURL, userID))
		waitgroup.Done()
	}()
	go func() {
		balanceResponse, _ = http.Get(fmt.Sprintf("%s/balance/%s", serverURL, userID))
		waitgroup.Done()
	}()
	go func() {
		debtsResponse, _ = http.Get(fmt.Sprintf("%s/user-debts/%s", serverURL, userID))
		waitgroup.Done()
	}()
	waitgroup.Wait()

	var userInfo map[string]string
	unmarshalResponse(userResponse, &userInfo)
	var userBalance map[string]string
	unmarshalResponse(balanceResponse, &userBalance)
	var userDebts []map[string]string
	unmarshalResponse(debtsResponse, &userDebts)

	return UserStatus{
		ID:            userInfo["id"],
		Name:          userInfo["name"],
		BalanceAmount: userBalance["amount"],
		Debts:         userDebts,
	}, nil
}
```

3. Una llamada utilizando funciones anónimas, closures, goroutines y un waitgroup para esperar a las respuesta de los 3 endpoints.

En este caso:
* Declaramos 3 channels `userResponse, balanceResponse, debtsResponse`.
* Cada request http lo hacemos dentro de una lambda y cada lambda se ejecuta en una goroutine diferente y una vez que tiene la respuesta se envía por el channel a la goroutine principal.
* Para obtener cada response, lo escuchamos del channel correspondiente, ejemplo: `<-userResponse`
* Se unen las respuesta como en el caso anterior.

```go
func GetUserStatusAsyncChannels(serverURL, userID string) (UserStatus, error) {

	userResponse := make(chan *http.Response)
	balanceResponse := make(chan *http.Response)
	debtsResponse := make(chan *http.Response)
	defer close(userResponse)
	defer close(balanceResponse)
	defer close(debtsResponse)

	go func() {
		result, _ := http.Get(fmt.Sprintf("%s/users/%s", serverURL, userID))
		userResponse <- result
	}()
	go func() {
		result, _ := http.Get(fmt.Sprintf("%s/balance/%s", serverURL, userID))
		balanceResponse <- result
	}()
	go func() {
		result, _ := http.Get(fmt.Sprintf("%s/user-debts/%s", serverURL, userID))
		debtsResponse <- result
	}()

	var userInfo map[string]string
	unmarshalResponse(<-userResponse, &userInfo)
	var userBalance map[string]string
	unmarshalResponse(<-balanceResponse, &userBalance)
	var userDebts []map[string]string
	unmarshalResponse(<-debtsResponse, &userDebts)
	return UserStatus{
		ID:            userInfo["id"],
		Name:          userInfo["name"],
		BalanceAmount: userBalance["amount"],
		Debts:         userDebts,
	}, nil
}
```

Para comprobar hacemos un test que levanta un servidor http y ejecuta los 3 clientes y podemos observar los resultados.

```go
func TestGetUserStatus(t *testing.T) {
	srv := httptest.NewServer(handler())
	defer srv.Close()
	userID := "2"

	start := time.Now()
	result, _ := GetUserStatusSync(srv.URL, userID)
	elapsed := time.Since(start)
	log.Printf("GetUserStatusSync took %s\n", elapsed)

	start = time.Now()
	result, _ = GetUserStatusAsyncWaitGroup(srv.URL, userID)
	elapsed = time.Since(start)
	log.Printf("GetUserStatusAsyncWaitGroup took %s\n", elapsed)

	start = time.Now()
	result, _ = GetUserStatusAsyncChannels(srv.URL, userID)
	elapsed = time.Since(start)
	log.Printf("GetUserStatusAsyncChannels took %s\n", elapsed)
	...
}
```

Resultados:
* Caso bloqueante tarda 750 ms, que es el total sumado de los delays.
* Casos 2 y 3 tardan 350 ms aprox, que es el mayor de los delays.

```sh
2019/10/12 17:21:35 server listening connections
2019/10/12 17:21:36 GetUserStatusSync took 763.135932ms
2019/10/12 17:21:36 GetUserStatusAsyncWaitGroup took 352.60757ms
2019/10/12 17:21:36 GetUserStatusAsyncChannels took 355.580379ms
PASS
```
El código completo y los tests aplicados se puede encontrar [aqui](https://github.com/jegutierrez/functional_patterns_go/tree/master/http)

Aprovechando la capacidad de pasar funciones podemos hacer patrones de concurrencia muy elegantes, manteniendo la simpleza en nuestro código.

### Conclusiones