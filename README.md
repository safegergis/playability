# Playability

Playability is a full stack web app dedicated to making video games more accessible and inclusive for everyone, especially those with disabilities. Our mission is to empower gamers of all abilities by providing comprehensive accessibility information and fostering a community-driven approach to game accessibility.

## Mission Statement

At Playability, we believe that everyone should have the opportunity to enjoy video games, regardless of their physical or cognitive abilities. We strive to break down barriers in gaming by providing accurate, user-driven accessibility information and advocating for more inclusive game design. Our goal is to create a world where no gamer is left behind due to lack of accessibility features.

## Images
<img width="1509" alt="1" src="https://github.com/user-attachments/assets/53403101-f4bd-42d2-898f-a7ed39e86505">
<img width="1509" alt="2" src="https://github.com/user-attachments/assets/d08bfbbf-42cd-4ab4-b2c9-dab104d785d6">
<img width="1509" alt="4" src="https://github.com/user-attachments/assets/cf2e78bf-6ea1-4869-9c77-dc2d3bc41690">

## Features

### User-Facing Features

- **User Registration and Authentication:** Secure user sign-up and login with email verification and password recovery.
- **Comprehensive Game Database:** Access a vast library of games with detailed information sourced from IGDB and PCGamingWiki.
- **Detailed Game Pages:** View comprehensive game information including:
  - Cover art and screenshots
  - Supported platforms
  - Release date and developer information
  - Genre and tags
  - Accessibility features with user-submitted ratings and comments
- **User-Driven Accessibility Feedback:**
  - Submit detailed reports on a game's accessibility features
  - Rate the effectiveness of existing accessibility options
  - Provide comments and tips for other users with similar needs
- **Accessibility Information:** Detailed breakdown of features such as:
  - Closed captions and subtitle options
  - Colorblind modes and visual assistance settings
  - Full controller support and button remapping capabilities
  - Difficulty settings and assist modes
  - Text-to-speech and screen reader compatibility
- **Accessibility Scoring System:** Aggregate user feedback to provide overall accessibility scores for games across different categories
- **Responsive Design:** Mobile-friendly interface using Tailwind CSS, ensuring the platform is accessible across all devices

### Technical Features

- **RESTful API:** Well-structured API endpoints for seamless communication between frontend and backend.
- **JWT Authentication:** Secure, token-based authentication system for protected routes and user sessions.
- **Database Caching:** Implemented Redis caching to improve performance and reduce load on the primary database.
- **CI/CD Pipeline:** Automated build, test, and deployment processes using GitHub Actions.
- **Microservices Architecture:** Modular backend services for improved scalability and maintainability.
- **Data Analytics:** Integration with analytics tools to track user behavior and improve the platform based on usage patterns.
- **Accessibility-First Design:** Ensure that the Playability platform itself adheres to WCAG guidelines and is fully accessible to users with disabilities.

## Technologies

- **Frontend:**
  - [Nuxt 3](https://nuxt.com/) - Vue.js framework for server-side rendering and improved SEO
  - [Vue.js](https://vuejs.org/) - Progressive JavaScript framework for building user interfaces
  - [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework for rapid UI development
  - [ShadCN Vue](https://shadcn.com/) - Accessible and customizable UI component library
  - [Vee Validate](https://vee-validate.logaretm.com/v4/) - Form validation library for Vue.js
  - [Vue Carousel](https://ismail9k.github.io/vue3-carousel/) - Accessible carousel component for featured games
- **Backend:**
  - [Go](https://golang.org/) - High-performance programming language for backend services
  - [Chi Mux](https://github.com/go-chi/chi) - Lightweight and flexible HTTP router for Go
  - [PostgreSQL](https://www.postgresql.org/) - Robust relational database for data persistence
  - [IGDB API](https://api-docs.igdb.com/) - Comprehensive source for video game data
  - [PCGamingWiki API](https://pcgamingwiki.com/api.php) - Additional source for PC game accessibility information
  - [JWT](https://jwt.io/) - Secure user authentication and authorization
