# Hash Verifier
* O arquivo utilizado foi o [rockyou.txt](https://github.com/praetorian-code/Hob0Rules/blob/master/wordlists/rockyou.txt.gz)
* O servidor possui um script que divide os arquivos utilizando shell que funciona somente no linux

## Servidor

Vá até a pasta do servidor e digite o seguinte comando no terminal:

``` go run *.go ```

Após isso o servidor estará rodando na porta `8080`.

Para verificar alguma hash, basta passar a hash como parâmetro e um arquivo `txt` em uma request do tipo `POST` 
e `multipart/form-data` com o nome `wordlist`. Para fins de verificação, o servidor possui uma rota para 
testar a hash de maneira síncrona (`localhost:8080/sync`) e assíncrona (`localhost:8080/async`). 

Por exemplo:

``` CURL http://localhost:8080/async/fd9cabd4def5137a73d682f4dd963e57 ```
