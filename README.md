# SDKs Tábua da Maré

Este repositório reúne diversas SDKs para integração com a API de Tábua da Maré, facilitando o acesso a dados de marés para diferentes linguagens de programação.

## Sobre o Projeto

A Tábua da Maré é um projeto que disponibiliza informações sobre marés através de uma API REST. Este repositório centraliza as SDKs desenvolvidas para consumir essa API em diferentes linguagens, tornando a integração mais simples e acessível para desenvolvedores.

Para mais informações sobre o serviço, acesse: [tabuamare.devtu.qzz.io](https://tabuamare.devtu.qzz.io)

## Projetos Relacionados

Este projeto faz parte de um ecossistema maior:

- **API Principal**: [tabua_mare_api](https://github.com/Ddiidev/tabua_mare_api) - A API REST que serve os dados de tábua da maré
- **Processamento de Dados**: [tabua_mare_convert_pdf2db](https://github.com/Ddiidev/tabua_mare_convert_pdf2db) - Projeto responsável por processar os PDFs de tábua da maré e converter para o banco de dados que alimenta a API

## SDKs Disponíveis

### JavaScript/TypeScript

SDK para Node.js e navegadores com suporte completo a TypeScript.

- 📦 **Instalação**: `npm install tabua-mare-sdk`
- 📖 **Documentação**: [README completo](./JsTs/README.md)
- ✨ **Características**: Zero dependências, totalmente tipado, funciona em Node.js e Browser

### Go

SDK oficial em Go com suporte a context e configuração flexível.

- 📦 **Instalação**: `go get github.com/Ddiidev/sdks-tabua-mare/go`
- 📖 **Documentação**: [README completo](./go/README.md)
- ✨ **Características**: Zero dependências externas, validação de parâmetros, tratamento robusto de erros

## Início Rápido

Escolha a SDK da sua linguagem preferida e siga a documentação específica:

```bash
# JavaScript/TypeScript
npm install tabua-mare-sdk

# Go
go get github.com/Ddiidev/sdks-tabua-mare/go
```

Cada SDK possui sua própria documentação com exemplos de uso, instalação e configuração. Consulte os links acima para mais detalhes.

## Como Contribuir

Contribuições são bem-vindas! Você pode contribuir de várias formas:

### 1. Melhorando SDKs Existentes

- Faça um fork do repositório
- Crie uma branch para sua feature (`git checkout -b feature/minha-feature`)
- Faça commit das suas alterações (`git commit -m 'Adiciona nova feature'`)
- Faça push para a branch (`git push origin feature/minha-feature`)
- Abra um Pull Request

### 2. Criando uma Nova SDK

Se você deseja criar uma SDK para uma nova linguagem:

1. Crie uma pasta na raiz do projeto com o nome da linguagem (ex: `python`, `ruby`, `rust`)
2. Estruture o projeto seguindo as convenções da linguagem
3. Inclua um README específico com:
   - Instruções de instalação
   - Exemplos de uso
   - Documentação da API
   - Guia de testes
4. Adicione testes unitários e de integração
5. Atualize este README principal incluindo a nova SDK na seção "SDKs Disponíveis" com:
   - Nome da linguagem
   - Comando de instalação
   - Link para o README específico
   - Principais características

### 3. Reportando Bugs ou Sugerindo Melhorias

- Abra uma [issue](https://github.com/Ddiidev/sdks-tabua-marea/issues) descrevendo o problema ou sugestão
- Seja claro e forneça exemplos quando possível

### 4. Diretrizes Gerais

- Mantenha o código limpo e bem documentado
- Siga as convenções de código da linguagem específica
- Adicione testes para novas funcionalidades
- Atualize a documentação quando necessário
- Seja respeitoso e construtivo nas discussões

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## Contato

Para dúvidas ou sugestões sobre o projeto, abra uma issue ou entre em contato através do repositório principal.
