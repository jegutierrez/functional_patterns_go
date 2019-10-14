# Functional patterns en Go
Recopilación de algunos patrones útiles que he identificado en el dia a dia trabajando con Go.

**Nota:**
Este repo NO es una introducción a Go, se asume que el lector tiene idea de la sintaxis básica del leguaje, manejo de funciones y conocimiento básico de algunos paquetes de la standard library como `net/http`. En caso de que no sea así, se recomienda primero dar una mirada al [tour de Go](https://tour.golang.org/welcome/1) que es una guia oficial y muy completa al lenguaje.

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

### Funciones como valor

En Go se acostumbra utilizar funciones como valor y podemos ver que es muy utilizado en la standard library del lenguaje.

Un ejemplo de esto es el paquete `net/http`, ampliamente utilizado en el ecosistema de Go para hacer servicios http (APIs REST por ejemplo).

Para ver el uso de funciones primero vamos a describir como acostumbramos escribir un servidor HTTP básico usando el paquete `net/http`.

Un concepto fundamental en los servidores `net/http` son los handlers. Un handler es un objeto que implementa la interfaz `http.Handler`. Una forma común de escribir un handler es mediante el uso del adaptador `http.HandlerFunc` que es tipo que describe la firma adecuada.

Las funciones que sirven como handlers toman un `http.ResponseWriter` y un `http.Request` como argumentos. El writer de respuestas se utiliza para completar la respuesta HTTP.

Cuando creamos una aplicación web, probablemente haya alguna funcionalidad compartida que queremos ejecutar para muchos (o incluso todos) los request HTTP. Es posible que deseemos loggear cada solicitud, comprimir cada respuesta, hacer validaciones o actualizar un caché antes de realizar un procesamiento pesado.

Una forma de organizar esta funcionalidad compartida es configurarla como middleware. Que es un código autónomo que actúa de forma independiente con cada request, antes o después de los handlers de aplicaciones normales. En Go, un lugar común para usar middleware es entre un servidor y sus handlers.

Los middlewares http son nada más que funciones que reciben y retornan un `http.HandlerFunc` y dentro podemos hacer operaciones necesarias sobre el request y/o response. 

Veamos un ejemplo: este middleware valida que un usuario esté autenticado, sino retorna un "404 Not found" y evita que se obtenga la información.

```go
func onlyAuthenticated(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAuth(r) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h(w, r)
	}
}
```

Luego al llamar a nuestro handler http lo wrappeamos con el middleware de la siguiente manera.

```go
srv := http.NewServeMux()
srv.HandleFunc("/balance/", onlyAuthenticated(balanceHandler))
```

Otro ejemplo que podemos dar es un validador, segun el tipo de datos guardamos en un map la funcion de validacion correspondiente

Supongamos que tenemos los tipos Movement y validator

```go
type Movement struct {
	ID           int
	Amount       float64
	Fee          float64
	MovementType string
}

type validator func(Movement) bool
```

Luego tenemos un mapa de funciones de validacion. Es importante destacar que Go trata las funciones como valor y por eso las podemos guardar en un mapa.

En este caso definimos que tenemos dos tipos de datos income y expense
* El income debe ser positivo y el movimiento debe tener un fee asociado
* El expense debe ser negativo

```go
var MovementValidator = map[string]validator{
	"income": func(m Movement) bool {
		if m.Amount <= 0 {
			return false
		}
		if m.Fee <= 0 {
			return false
		}
		return true
	},
	"expense": func(m Movement) bool {
		if m.Amount >= 0 {
			return false
		}
		return true
	},
}
```

Luego podemos usar nuestro validador de la siguiente manera:

* Declaramos 2 movimientos validos y 1 invalido
* Imprimimos el ID del movimiento si es invalido
```go
func main() {
	validIncome := Movement{
		ID:           1,
		Amount:       10,
		Fee:          1,
		MovementType: "income",
	}
	validExpense := Movement{
		ID:           2,
		Amount:       -10,
		MovementType: "expense",
	}
	invalidIncomeMov := Movement{
		ID:           3,
		Amount:       10,
		MovementType: "income",
	}

	if !MovementValidator[validIncome.MovementType](validIncome) {
		log.Printf("Invalid movement %d", validIncome.ID)
	}
	if !MovementValidator[validExpense.MovementType](validExpense) {
		log.Printf("Invalid movement %d", validExpense.ID)
	}
	if !MovementValidator[invalidIncomeMov.MovementType](invalidIncomeMov) {
		log.Printf("Invalid movement %d", invalidIncomeMov.ID)
	}
}
```

### Closures

Un closure es la combinación de una función y el ámbito en el que se declaró dicha función. Y la particularidad es que la función definida en el closure "recuerda" el entorno en el que se ha creado y puede acceder a valores o punteros de ese entorno en cualquier momento.

Esto es una de las cosas mas poderosas para tener en cuenta al trabajar con funciones en Go.

Veamos 3 casos que he visto implementados en programas productivos, donde los closures puede ser muy útiles:

1. Filtrado de datos Generico

Supongamos que tenemos tipos que comparten un tipo de dato en común o ninguno, pero tenemos repetida la logica de filtrados en varias partes de nuestro programa.

Podemos aprovechar el pase de funciones y los closures para hacer una funcion generica.

Tenemos los tipos AccountMovement y Debt, que no tienen mucha relacion entre ellos.
```go
type AccountMovement struct {
	ID     int
	From   string
	To     string
	Amount float64
}

type Debt struct {
	ID     int
	UserID int
	Reason string
	Amount float64
}
```

Y luego tenemos un slice de AccountMovements y Debts.
```go
movements := []AccountMovement{
	{ID: 1, From: "a", To: "b", Amount: 7},
	{ID: 2, From: "c", To: "b", Amount: 14},
	...
}
debts := []Debt{
	{ID: 1, Reason: "x", UserID: 4, Amount: 16},
	{ID: 2, Reason: "x", UserID: 2, Amount: 4},
	...
}
```

En algunos lenguales como java o javascript tenemos la funcion helper `filter()` que puede ser usadad sobre una lista, pasando una funcion predicado nos retorna una nueva lista de los elemtos donde ese predicado es `true`, en Go no tenemos estas funciones en la standar library, pero podemos contruirlas con un poco de ayuda de los closures.

```go
func Filter(l int, predicate func(int) bool, appender func(int)) {
	for i := 0; i < l; i++ {
		if predicate(i) {
			appender(i)
		}
	}
}
```
Y para llamarla en las listas de tipo AccountMovement y Debt hacemos lo siguiente:

* Declaramos un slice de elementos AccountMovement fuera de la llamada a filter
```go
var bigMovements []AccountMovement
Filter(len(movements), func(i int) bool {
	return movements[i].Amount > 20
}, func(i int) {
	bigMovements = append(bigMovements, movements[i])
})

var bigDebts []Debt
Filter(len(debts), func(i int) bool {
	return debts[i].Amount > 20
}, func(i int) {
	bigDebts = append(bigDebts, debts[i])
})
```

2. Testing

Otro luegar donde es muy común utilizar tecnicas funcionales es en los tests.

Pasando closures a funciones en los tests podemos mockear funciones que nos interesa testear, ademas de poder hacer asserts dentro del closure conservando el scope de cada ejecución de un tests.

3. Handler http


### Partial Aplication

1. Testing

Para los tests utilizamos esta tecnica, de manera de aplicar el context de la ejecución de un test.

2. Http handler



### Concurrencia y funciones

Frecuentemente utilizamos técnicas funcionales para trabajar con concurrencia en Go como funciones lambda o closures. 

Vamos a mostrar el ejemplo muy común de hacer llamadas a múltiples APIs y luego unir los resultados, lo que buscamos es, hacer múltiples requests http en paralelo y luego esperar a que vuelvan las respuestas para procesarlas.

En Go no tenemos en la standard library algo parecido a las promesas en javascript o Futures en java, sin embargo utilizado goroutines, waitgroups y channel podemos construir nuestra propia manera de resolver promesas al estilo simple de Go.

Primero mostramos un ejemplo de 3 endpoint, con delays:

* /users 		-> 150 ms
* /balance 		-> 350 ms
* /user-debts 	-> 250 ms


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