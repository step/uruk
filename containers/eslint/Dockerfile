FROM cytopia/eslint
WORKDIR /source
ENTRYPOINT [ "eslint", ".", "--no-eslintrc" ,"-f", "json", "-o", "/results/result.json"]