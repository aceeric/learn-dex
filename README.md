# Dex Learning

This is a project to learn [Dex IdP](https://www.google.com/url?sa=t&source=web&rct=j&opi=89978449&url=https://dexidp.io/&ved=2ahUKEwiO9_Dbt6GHAxVqEVkFHcFQAiwQFnoECBcQAQ&usg=AOvVaw2kFMtDZIG3uN8LlA284iYF), and integrate Dex OIDC authentication and authorization into a Kubernetes workload.

This project is organized into a series of Labs. The labs are structured such that others should be able to repeat the steps using copy/paste. All labs were performed on a Ubuntu 22.04 Desktop, but any Linux environment with Bash shell support should work fine.

## Labs

1. **[First Touch](first-touch.md)**: In this lab, we git clone the Dex GitHub repository, then run the Dex server with the built-in Mock OIDC Provider as a Docker container and access the server using the example client from the cloned repo.
1. **[Configure an external OIDC Provider](mock-oidc-provider.md)**: In this lab, we build a different Mock OIDC Provider and integrate it into Dex as a new connector.

## To Do Labs

- Run a Python web app on the desktop that integrates with Dex
- Run Dex in Kubernetes and run the example client app in k8s and show the integration
- Build and run a trivial Python app in k8s w/ Dex integration
- Integrate Dex with Nginx in k8s and test using the trivial Python app
