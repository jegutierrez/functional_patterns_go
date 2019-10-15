# Functional patterns en Go

Recopilación de algunos patrones útiles que he identificado en el dia a dia trabajando con Go.

**Nota:**
Este repo **NO** es una introducción a Go, se asume que el lector tiene idea de la sintaxis básica del lenguaje, manejo de funciones y conocimiento básico de algunos paquetes de la standard library como `net/http` y `testing`. En caso de que no sea así, se recomienda primero dar una mirada al [tour de Go](https://tour.golang.org/welcome/1) que es una guía oficial y muy completa al lenguaje.

## Entonces, ¿programación funcional en Go?

Go tiene la capacidad de usar funciones como ciudadanos de primera clase, es decir, se pueden utilizar las funciones para pasarlas como parámetro o retornarlas como valor de otra función (esto es conocido como funciones de orden superior), sin embargo, como todo en software, es acerca de tradeoffs, lo recomendado es siempre preferir la claridad del código antes que cualquier patrón o paradigma.

En el siguiente post vamos a describir una serie de técnicas de programación funcional, aprovechando esta capacidad de trabajar con funciones como valor.

## Que implica hacer programación funcional

La programación funcional nos trae dos grandes restricciones:

* Usar funciones para todo: hacer todas las operaciones con funciones.
* No mutar estado: no mutar valores una vez declarados, no tener estructuras de datos mutables, no tener side effects dentro de nuestras funciones.

Algunos de los conceptos clave en programación funcional son los siguientes:

## Recursividad

Cuando hablamos de programación funcional (FP) tenemos dos requisitos importantes, trabajar con funciones y no mutar estado. Y una de las primeras técnicas que nos viene a la cabeza es la recursividad, que es la capacidad de una función de llamarse a sí misma. Esto es necesario porque es la manera de hacer loops en FP.

Por lo tanto lo primero que no se puede usar para mantener inmutabilidad en nuestra aplicación es **for loops**, ya que vamos mutando una variable que toma diferente valor en cada iteración.

Otra de las cosas que acostumbramos a usar en Go son estructuras de datos a las que vamos agregando elementos, esto tampoco podríamos si seguimos las reglas de FP, aunque en el caso del **slice** usando la función `append()` no estamos violando la regla, ya que nos retorna un nuevo slice cada vez que se llama.

#### ¿Es buena idea no usar for loops en Go?

Veamos con un ejemplo y pruebas de performance:

Escribir una función que cuente cuántas maneras posibles de dar cambio hay para un monto dado con una lista de denominaciones de monedas.

Por ejemplo, hay 3 maneras de dar cambio para \$4 con monedas de \$1 y \$2
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

## Funciones como valor

En Go se acostumbra utilizar funciones como valor y es ampliamente utilizado en la standard library del lenguaje.

Un ejemplo de esto es el paquete `net/http`, ampliamente utilizado en el ecosistema de Go para hacer servicios http (APIs REST por ejemplo).

Para ver el uso de funciones primero vamos a describir como acostumbramos escribir un servidor HTTP básico usando el paquete `net/http`.

Un concepto fundamental en los servidores `net/http` son los handlers. Un handler es un objeto que implementa la interfaz `http.Handler`. Una forma común de escribir un handler es mediante el uso del adaptador `http.HandlerFunc` que es tipo que describe la firma adecuada.

Las funciones que sirven como handlers toman un `http.ResponseWriter` y un `http.Request` como argumentos. El writer de respuestas se utiliza para completar la respuesta HTTP.

Cuando creamos una aplicación web, probablemente haya alguna funcionalidad compartida que queremos ejecutar para muchos (o incluso todos) los request HTTP. Es posible que deseemos loggear cada solicitud, comprimir cada respuesta, hacer validaciones o actualizar un caché antes de realizar un procesamiento pesado.

Una forma de organizar esta funcionalidad compartida es configurarla como middleware. Que es un código autónomo que actúa de forma independiente con cada request, antes o después de los handlers de aplicaciones normales. En Go, un lugar común para usar middleware es entre un servidor y sus handlers.

Los middlewares http son nada más que funciones que reciben y retornan un `http.HandlerFunc` y dentro se pueden hacer operaciones necesarias sobre el request y/o response. 

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

Otro ejemplo es un validador, según el tipo de datos guardamos en un map la función de validación correspondiente

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

Luego tenemos un mapa de funciones de validación. Es importante destacar que Go trata las funciones como valor y por eso las podemos guardar en un Map.

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

Luego se usa el validador de la siguiente manera:

* Declaramos 2 movimientos válidos y 1 inválido
* Imprimimos el ID del movimiento si es inválido
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

## Closures y partial application

Un closure es la combinación de una función y el ámbito en el que se declaró dicha función. Y la particularidad es que la función definida en el closure "recuerda" el entorno en el que se ha creado y puede acceder a valores o punteros de ese entorno en cualquier momento.

Esto es una de las cosas más poderosas para tener en cuenta al trabajar con funciones en Go.

Veamos 3 casos que he visto implementados en programas productivos, donde los closures puede ser muy útiles:

**1. Filtrado de datos Genérico**

Supongamos que tenemos tipos que comparten un tipo de dato en común o ninguno, pero tenemos repetida la lógica de filtrados en varias partes de nuestro programa.

Aprovechando el pase de funciones y los closures para hacer una función genérica.

Tenemos los tipos AccountMovement & Debt, que no tienen mucha relación entre ellos.
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

En algunos lenguajes como java o javascript tenemos la función helper `filter()` que puede ser usada sobre una lista, pasando una función predicado nos retorna una nueva lista de los elementos donde ese predicado es `true`, en Go no tenemos estas funciones en la standard library, pero podemos construirlas con un poco de ayuda de los closures.

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

**2. Testing**

Otro lugar donde es muy común utilizar técnicas funcionales es en los tests.

Veamos cómo se ven los test en Go utilizando el paquete `testing`:

```go
import "testing"

func TestAbc(t *testing.T) {
    t.Error() // para indicar que el test falló
}
```

Veamos un caso real:

Tenemos una interfaz DB, con un método para guardar un usuario en la base de datos y el tipo MySQL que lo implementa.

```go
type DB interface {
	SaveUser(u User)
}

type MySQL struct{}

func (m MySQL) SaveUser(u User) {
	// DB save
}
```

Pasando closures a funciones en los tests se pueden mockear funciones que nos interesa testear, además de poder hacer asserts dentro del closure conservando el scope de cada ejecución de un tests.

Definimos un tipo MySQL mock:

```go
type MockDB struct {
	MockSaveUserFn func(User)
}

func (m MockDB) SaveUser(u User) {
	m.MockSaveUserFn(u)
}
```

Y un helper que recibe el context del test `*testing.T` y en retorna la función para guardar un usuario que utiliza el `t` declarado en el closure.

```go
func helperMockDB(t *testing.T) func(User) {
	t.Helper()

	return func(u User) {
		if u.ID != 0 {
			t.Errorf("user ID must not be preset")
		}
	}
}
```

Dentro de un test creamos el handler, inyectando el contexto del test y luego se ejecuta el assert dentro de la función SaveUser.

```go
func TestHttpHandler(t *testing.T) {
	body := strings.NewReader(`{"name": "john"}`)
	req, err := http.NewRequest("POST", "/users", body)
	if err != nil {
		t.Fatal(err.Error())
	}

	res := httptest.NewRecorder()

	saveFn := helperMockDB(t)
	mockDB := MockDB{
		MockSaveUserFn: saveFn,
	}

	handler := saveUserHandler(mockDB)

	handler(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
```

**3. Handler http**

Para hacer servidores http uno de los componentes clave es el http handler.

Veamos un caso real, suponiendo que tenemos una dependencia como Newrelic para hacer tracing de requests en nuestra API.

```go
type FakeNewrelic struct {
	Name string
}

func NewRelicTracer(name string) FakeNewrelic {
	return FakeNewrelic{
		Name: fmt.Sprintf("trace %s", name),
	}
}

func (n *FakeNewrelic) Trace() {
	log.Printf(n.Name)
}
```

Algo que resulta muy útil es, en vez de que nuestro handler sea un `http.HandlerFunc`, que sea una función que recibe los parámetros necesarios y retorna un `http.HandlerFunc`, esto nos permite recibir parámetro y crear un entorno closure donde se puede inicializar funcionalidad antes de crear nuestro handler en sí. Para que quede más claro, veamos un ejemplo.

Después de tener definida nuestra dependencia (Newrelic) vamos a ver como utilizarla en nuestro handler:

* La función `userHandler` recibe un delay para utilizar dentro del handler.
* `userHandler` es un closure que nos permite declarar e inicializar cualquier dependencia antes de retornar el handler. En nuestro caso inicializamos un tracer y declaramos un tipo response.
* Despues dentro del handler `func(w http.ResponseWriter, r *http.Request)` podemos utilizar el tracer `nr.Trace()` y el delay que recibe como parámetro el `userHandler` de la siguiente manera `time.Sleep(delayMs * time.Millisecond)`.

```go
func userHandler(delayMs time.Duration) http.HandlerFunc {

	nr := NewRelicTracer("users")

	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		nr.Trace()

		userID := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.Println("userID is not a number")
			w.WriteHeader(400)
		}
		user := response{ID: id, Name: "user" + userID}

		time.Sleep(delayMs * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
```

## Concurrencia

Frecuentemente utilizamos técnicas funcionales para trabajar con concurrencia en Go como funciones lambda o closures. 

Vamos a mostrar el ejemplo muy común de hacer llamadas a múltiples APIs y luego unir los resultados, lo que buscamos es, hacer múltiples requests http en paralelo y luego esperar a que vuelvan las respuestas para procesarlas.

En Go no tenemos en la standard library algo parecido a las promesas en javascript o Futures en java, sin embargo utilizado goroutines, waitgroups y channel podemos construir nuestra propia manera de resolver promesas al estilo simple de Go.

Primero mostramos un ejemplo de 3 endpoint, con delays:

* /users 		-> 150 ms
* /balance 		-> 350 ms
* /user-debts	-> 250 ms


**1. Cliente http bloqueante**
Primero hacemos un cliente con los 3 request bloqueantes, sin paralelismo:

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

**2. Cliente asíncrono con waitgroups**
Una llamada utilizando funciones anónimas, closures, goroutines y un waitgroup para esperar a las respuesta de los 3 endpoints.

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

**3. Cliente asíncrono con channels**
Una llamada utilizando funciones anónimas, closures, goroutines y un waitgroup para esperar a las respuesta de los 3 endpoints.

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

Para comprobar que se ejecutan en paralelo escribimos un test que levanta un servidor http y ejecuta los 3 clientes y observar los resultados.

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

Aprovechando la capacidad de pasar funciones se pueden hacer patrones de concurrencia muy elegantes, manteniendo la simpleza en nuestro código.

## Conclusiones

* Utilizando funciones como ciudadanos de primera clase en Go se pueden construir aplicaciones flexibles y sin dejar de ser idiomáticas para el ecosistema.
* Closures son una de las herramientas más poderosas que tenemos en Go, con los cuales se puede construir funciones genéricas, inicializar dependencias en los handlers, mockear dependencias en los test y construir patrones de concurrencia.
* No conviene hacer programación funcional pura en Go, ya que no es idiomático en Go, la sintaxis no es amigable para trabajar con funciones. En todo caso, lo ideal es priorizar la claridad del código para el equipo que lo mantiene.
* En muchos de los casos es menos eficiente la versión funcional.
* Muchas veces usamos patrones sin darnos cuenta, es bueno identificarlos y utilizarlos en nuestras aplicaciones.
