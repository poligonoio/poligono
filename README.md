<!-- ABOUT THE PROJECT -->
# Poligono: The AI-Powered Data Analytics Platform

Poligono is an open-source data analytics platform that empowers users to effortlessly explore, visualize, and gain insights from their data using natural language (No SQL expertise required!).

<!-- KEY FEATURES -->
## ✨ Key Features

* **Natural Language Interface:** Ask questions about your data in plain English.
* **AI-Powered Insights:** Discover hidden trends and patterns.
* **Integrations:** Connect to your favorite data sources.
* **Extensible:** Build custom modules and plugins.

<!-- GETTING STARTED -->
## 🚀 Getting Started

Poligono is designed to run in a containerized environment. Here's the simplest way to get started:

1. **Prerequisites:**
   * Docker: [Install Docker](https://www.docker.com/products/docker-desktop).
   * API Keys: Obtain the required API keys for Gemini and Infisical.
2. **Run with Docker:**

    ```bash
    docker run -it \
        --name poligono \
        -e GEMINI_API_KEY=your_gemini_api_key \
        -e MONGODB_URI=your_mongodb_uri \
        -e INFISICAL_CLIENT_ID=your_infisical_client_id \
        -e INFISICAL_CLIENT_SECRET=your_infisical_client_secret \
        -e INFISICAL_PROJECT_ID=your_infisical_project_id \
        ghcr.io/poligonoxyz/poligono:latest
    ```

Once the Docker installation is complete, go to [http://localhost:8888/v1/swagger/index.html#/](http://localhost:8888/v1/swagger/index.html#/) to access the Poligono OpenAPI Specification from your browser.

<!-- CONTRIBUTING -->
## Contributing

All code contributions, including those of people having commit access, must go through a pull request and be approved by a core developer before being merged. This is to ensure a proper review of all the code.

We truly ❤️ pull requests! If you wish to help, you can learn more about how you can contribute to this project in the [contribution guide](CONTRIBUTING.md).
