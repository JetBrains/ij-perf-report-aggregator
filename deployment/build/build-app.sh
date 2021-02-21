set -e

cd dashboard-v2
node .yarn/releases/yarn-berry.cjs install --immutable
node .yarn/releases/yarn-berry.cjs build

ko resolve -f deployment/report-aggregator/values.yaml --tags "$BUILD_NUMBER" > deployment/report-aggregator/values-resolved.yaml
helm upgrade report-aggregator ./deployment/report-aggregator --install --cleanup-on-fail -f deployment/report-aggregator/values-resolved.yaml -f jb/values.yaml