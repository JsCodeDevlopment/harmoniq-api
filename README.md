# Harmoniq API - Core Service

Backend de alto desempenho desenvolvido em **Go** para suportar a infraestrutura do ecossistema Harmoniq. Esta API é responsável pela gestão de usuários, persistência de repertórios e processamento dinâmico de cifras.

<h1 align="center">
  <img alt="Banner do Harmoniq API" title="#Banner" style="object-fit: cover; width: 100%; max-height: 520px" src="public/preview.webp" />
</h1>

---

## 🚀 Funcionalidades do Ecossistema Backend

O Harmoniq API foi projetado para ser escalável e seguro, oferecendo respostas rápidas para que a música nunca pare.

### 1. 📡 Gestão de Cifras e Transposição

- **Motor de Transposição:** Lógica robusta para transpor nomes de acordes em tempo real, lidando com notas compostas e variações.
- **Normalização de Notas:** Sistema inteligente que trata sustenidos e bemóis de forma transparente.
- **Database de Músicas:** Integração com provedores externos para busca e recuperação de conteúdos.

### 2. 🔐 Autenticação e Segurança

- **Nível Elite (JWT):** Implementação segura de tokens para acesso às rotas privadas.
- **Controle de Sessão:** Gerenciamento eficiente de estados de usuário e permissões de acesso.
- **Proteção de Dados:** Criptografia avançada (bcrypt) para dados sensíveis e senhas.

### 3. 💾 Persistência de Repertórios (Setlists)

- **SQL Estruturado:** Uso de PostgreSQL (GORM) para armazenamento relacional de usuários, setlists e músicas.
- **Variações Customizadas:** Persistência de escolhas individuais de acordes por música no banco de dados.
- **Relacionamentos Complexos:** Gestão de acesso entre usuários e seus repertórios compartilhados.

### 4. ⚡ Desempenho e Cache

- **Integração com Redis:** Camada de cache para acelerar buscas frequentes e acelerar a entrega de conteúdo.
- **Processamento Paralelo:** Aproveitamento das rasteiras (goroutines) do Go para operações concorrentes.
- **CORS e Middleware:** Configuração profissional para integração segura com o frontend.

---

## 🛠️ Tecnologias

<div align="center">
  <img src="https://img.shields.io/badge/Go_1.25.0-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Gin_Gonic-008080?style=for-the-badge&logo=gin&logoColor=white" />
  <img src="https://img.shields.io/badge/GORM-3178C6?style=for-the-badge&logo=gorm&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/JWT_Authentication-000?style=for-the-badge&logo=json-web-tokens&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" />
  <img src="https://img.shields.io/badge/Vercel-000?style=for-the-badge&logo=vercel&logoColor=white" />
</div>

---

## 📅 Histórico de Versões

Confira todas as melhorias e novos recursos em nosso log oficial:
👉 **[Ver Histórico Completo em RELEASES.md](./RELEASES.md)**

---

## 👨💻 Time e Desenvolvimento

<div align="center">
  <img src="https://avatars.githubusercontent.com/u/100796752?s=400&u=ae99bd456c6b274cd934d85a374a44340140e222&v=4" width="100" style="border-radius: 50%" />
  <br>
  <strong>Jonatas Silva</strong>
  <br>
  Senior Software Engineer / CTO & Tech Lead at <a href="https://pokernetic.com/">PokerNetic</a>
</div>

---

## 📄 Licença

Este projeto é privado e de uso exclusivo da **Harmoniq Inc**.

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/JsCodeDevlopment">Jonatas Silva</a></sub>
</div>
