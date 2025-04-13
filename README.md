# ice-kube

`ice-kube` is a command-line tool designed to optimize infrastructure utilization and analyze resources within Kubernetes clusters. It offers functionalities such as scanning number of pods in a namespace or entire cluster and automatic deletion of pods in completed state based on predefined schedules along with dry run capabilities.

## Features

- **Scheduled Scaling**: Scale down pods attached to jobs based on predefined time schedules.
- **Resource Analysis**: Analyze resources in the cluster and determine dangling resources eligble for deleti.
- **Automated Operations**: Automatic delete operations to reduce manual intervention in cluster clean-up.
- **Dry Run Capability**: Preview delete action without executing it, ensuring safe operations.

## Installation

### Prerequisites

- Go (version 1.16 or higher)
- Access to a Kubernetes cluster
- `kubectl` configured to interact with your cluster

### Steps

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/cloud-sky-ops/ice-kube.git
   cd ice-kube
   ```

2. **Build the Application**:

   ```bash
   go build -o ice-kube main.go
   ```

3. **Verify Installation**:

   ```bash
   ./ice-kube --help
   ```

## Usage

### General Help

```bash
./ice-kube --help
```

### Commands

- **Scan**: Analyze current deployments and resource usage.

  ```bash
  ./ice-kube scan
  ```

- **Delete**: Remove deployments that meet specific criteria.

  ```bash
  ./ice-kube delete --namespace=default --label=app=example
  ```

- **Dry Run**: Preview actions without executing them.

  ```bash
  ./ice-kube delete --namespace=default --label=app=example --dry-run
  ```

### Flags

- `--cluster`: Name of the kubernetes cluster to set context.
- `--namespace`: Specify the Kubernetes namespace.
- `--dry-run`: Simulate the command without making changes.
- `--delete-before-hours`: Number of hours after which a completed Pod should be removed from the cluster.

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**: Click the "Fork" button at the top right of the repository page.
2. **Create a Branch**: Create a new branch for your feature or bugfix.
3. **Commit Changes**: Make your changes and commit them with clear messages.
4. **Push to Fork**: Push your changes to your forked repository.
5. **Submit a Pull Request**: Open a pull request to the main repository's `main` branch.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

Developed and maintained by [cloud-sky-ops](https://github.com/cloud-sky-ops).

## Sample executions

## Output for "ice-kube --help" 
![image](https://github.com/user-attachments/assets/d49a38ab-4094-4861-b366-5ab2a4993728)

## Output for delete command with --dry-run flag

![image](https://github.com/user-attachments/assets/60bf8306-472b-4ea3-83e5-20b8b696a8b0)

## Output for "delete command" without --dry-run flag
![image](https://github.com/user-attachments/assets/c9fe483b-c623-4a77-9578-5a0050cf5267)

## Output for "scan --help" command
![image](https://github.com/user-attachments/assets/522cdad0-615a-46b4-b43c-c0f395529392)

---
