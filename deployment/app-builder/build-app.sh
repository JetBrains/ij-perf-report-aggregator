set -e

pnpm --version
pnpm i --frozen-lockfile
pnpm lint
pnpm build

ko resolve -f deployment/report-aggregator/values.yaml --base-import-paths --tags "$BUILD_NUMBER" > deployment/report-aggregator/values-resolved.yaml
sed -i -e 's/:[a-z0-9]*@sha256/@sha256/g' deployment/report-aggregator/values-resolved.yaml
helm upgrade report-aggregator ./deployment/report-aggregator --install --cleanup-on-fail -f deployment/report-aggregator/values-resolved.yaml -f jb/values.yaml