set -e

cd dashboard
node .yarn/releases/yarn-sources.cjs install --immutable
node .yarn/releases/yarn-sources.cjs build

cd dashboard-v2
node .yarn/releases/yarn-sources.cjs install --immutable
node .yarn/releases/yarn-sources.cjs build

cd ..

ko resolve -f deployment/report-aggregator/values.yaml --tags "$BUILD_NUMBER" > deployment/report-aggregator/values-resolved.yaml
helm upgrade report-aggregator ./deployment/report-aggregator --install --cleanup-on-fail -f deployment/report-aggregator/values-resolved.yaml -f jb/values.yaml