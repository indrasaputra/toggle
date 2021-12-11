# Development Guide

- Fork or clone the project

    ```
    $ git@github.com:indrasaputra/toggle.git
    ```

- Create a meaningful branch

    ```
    $ git checkout -b <your-meaningful-branch>
    ```

    e.g:

    ```
    $ git checkout -b optimize-get-all-toggles
    ```

- Create some changes and their tests (unit test, integration test, and other test if any).

- If you want to generate files based on protocol buffer definition, run

    ```
    $ make gen.proto
    ```

    Alternatively, you can generate files using provided docker image.
    Run this command to use docker to generate files.

    ```
    $ make gen.proto.docker
    ```

- If you want to generate mock based on some interfaces, run

    ```
    $ make gen.mock
    ```

- Make sure you format/beautify the code by running

    ```
    $ make pretty
    ```

- As a reminder, always run the command above before add and commit changes.
    That command will be run in CI Pipeline to verify the format.

- Test your changes

    ```
    $ make test.unit
    ```

- Add, commit, and push the changes to repository

    ```
    $ git add .
    $ git commit -s -m "your meaningful commit message"
    $ git push origin <your-meaningful-branch>
    ```

    For writing commit message, please use [conventionalcommits](https://www.conventionalcommits.org/en/v1.0.0/) as a reference.

- Create a Pull Request (PR). In your PR's description, please explain the goal of the PR and its changes.

- Ask the other contributors to review.

- Once your PR is approved and its pipeline status is green, ask the owner to merge your PR.
