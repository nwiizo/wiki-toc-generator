# GitLab Wiki Table of Contents Generator

A command-line tool in Go that clones a GitLab Wiki repository and generates a table of contents in Markdown format.

## Installation

Clone the repository and build the tool:

```sh
git clone https://github.com/nwiizo/wiki-toc-generator.git
cd wiki-toc-generator
go build -o wiki-toc-generator main.go
```

## Usage

Generate the table of contents using the following command:

```sh
./wiki-toc-generator <wikiRepoURL>
```

Replace `<wikiRepoURL>` with the URL of your GitLab Wiki repository.

### Example

```sh
./wiki-toc-generator https://gitlab.com/yourusername/yourproject.wiki.git
```

Once the table of contents is generated, copy it to your Wiki page as needed.
