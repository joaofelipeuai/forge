# Forge Framework - Exemplo Avançado

Este exemplo demonstra todas as funcionalidades avançadas do Forge Framework.

## 🚀 Funcionalidades Demonstradas

### 🔐 Autenticação JWT
- Login com geração de token JWT
- Rotas protegidas que requerem autenticação
- Middleware de autenticação opcional e obrigatório
- Armazenamento seguro de tokens no localStorage

### 📁 Upload de Arquivos
- Upload múltiplo de arquivos
- Drag & drop interface
- Progress bar visual
- Validação de tipos de arquivo
- Armazenamento seguro de arquivos

### 🔄 WebSocket Real-time
- Conexão WebSocket bidirecional
- Chat em tempo real
- Sistema de broadcast
- Reconexão automática
- Status de conexão visual

### ⚡ Hot Reload
- Recarregamento automático durante desenvolvimento
- Monitoramento de mudanças em arquivos
- Atualização de templates em tempo real

### 🎨 Template Engine
- Sistema de templates dinâmico
- Funções personalizadas
- Recarregamento automático em modo desenvolvimento
- Integração com dados do servidor

### 🛡️ Middleware Integrado
- **CORS**: Configuração automática de headers
- **Rate Limiting**: Proteção contra spam
- **Recovery**: Recuperação de panics
- **Logger**: Log detalhado de requisições

## 📋 Como Executar

```bash
# Navegar para o diretório
cd examples/advanced

# Instalar dependências
go mod tidy

# Executar o servidor
go run advanced_server.go
```

O servidor iniciará na porta `:8080`

## 🌐 Endpoints Disponíveis

### Públicos
- `GET /` - Página inicial com informações
- `GET /template` - Demonstração do template engine
- `GET /health` - Health check do sistema

### Autenticação
- `POST /auth/login` - Login e geração de token JWT
- `GET /profile` - Perfil do usuário (requer JWT)

### Upload
- `POST /upload` - Upload múltiplo de arquivos
- `POST /upload/image` - Upload específico para imagens

### WebSocket
- `GET /ws` - Endpoint WebSocket
- `POST /broadcast` - Enviar mensagem broadcast

### Arquivos
- `GET /uploads/*` - Servir arquivos enviados

## 🧪 Como Testar

### 1. Autenticação JWT
```bash
# Fazer login
curl -X POST http://localhost:8080/auth/login

# Usar o token retornado
curl -H "Authorization: Bearer <token>" http://localhost:8080/profile
```

### 2. Upload de Arquivos
```bash
# Upload via curl
curl -X POST -F "files=@arquivo.txt" http://localhost:8080/upload
```

### 3. WebSocket
- Abra o navegador em `http://localhost:8080/template`
- Use a interface web para testar WebSocket

### 4. Interface Web Completa
Acesse `http://localhost:8080/template` para uma interface completa que demonstra todas as funcionalidades.

## 🎨 Interface Web

A interface web inclui:

- **Dashboard em tempo real** com estatísticas
- **Status indicators** para servidor, WebSocket e autenticação
- **Cards interativos** para cada funcionalidade
- **Chat em tempo real** via WebSocket
- **Upload com drag & drop**
- **Progress bars** e animações
- **Design responsivo** para mobile
- **Tema moderno** com gradientes e glassmorphism

## 🔧 Configurações

### JWT
```go
jwtConfig := forge.NewJWTConfig("your-secret-key")
jwtConfig.Expiration = 24 * time.Hour
```

### Upload
```go
uploadConfig := forge.NewUploadConfig("./uploads")
uploadConfig.MaxFileSize = 10 << 20 // 10MB
```

### Template Engine
```go
templateEngine := forge.NewTemplateEngine("templates", "html")
templateEngine.SetDevMode(true)
```

## 📁 Estrutura de Arquivos

```
examples/advanced/
├── advanced_server.go    # Servidor principal
├── test_jwt.go          # Teste específico de JWT
├── go.mod               # Dependências
├── templates/           # Templates HTML
│   └── index.html       # Interface principal
├── uploads/             # Arquivos enviados
└── README.md           # Esta documentação
```

## 🚀 Próximos Passos

1. **Personalizar** as configurações conforme sua necessidade
2. **Adicionar** mais rotas e funcionalidades
3. **Integrar** com banco de dados
4. **Implementar** autenticação real
5. **Adicionar** testes automatizados
6. **Deploy** em produção

## 🎯 Características Técnicas

- **Zero Dependencies**: Usa apenas a biblioteca padrão do Go
- **Type-Safe**: Aproveita o sistema de tipos do Go
- **High Performance**: Otimizado para alta performance
- **Production Ready**: Pronto para uso em produção
- **Extensível**: Fácil de estender e personalizar

---

**Forge Framework** - Construindo aplicações web modernas com Go! 🔨✨