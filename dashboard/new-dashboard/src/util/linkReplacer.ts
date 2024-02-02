export function replaceToLink(reason: string) {
  const ytRegexWithDescription = /https:\/\/youtrack.jetbrains.com\/issue\/(\w{2,}-\d{2,})\/(\S*)/
  const matchYT = reason.match(ytRegexWithDescription)
  let result = reason
  if (matchYT && matchYT.length > 2) {
    result = matchYT[1] + " " + matchYT[2].replaceAll("-", " ")
  }

  const style = 'class="underline decoration-dotted hover:no-underline"'

  result = result.replaceAll("https://youtrack.jetbrains.com/issue/", "")

  const issueRegex = /(\w{2,}-\d{2,})/g
  result = result.replaceAll(issueRegex, `<a ${style} href="https://youtrack.jetbrains.com/issue/$1">$1</a>`)

  const slack = /https:\/\/jetbrains.slack.com\/archives\/(\w){2,}\/p(\d{2,})/g
  const matchSlack = reason.match(slack)
  if (matchSlack) {
    result = result.replaceAll(slack, `<a ${style} href="${matchSlack[0]}">slack</a>`)
  }

  const commit = /https:\/\/jetbrains.team\/p\/(\w+)\/repositories\/(\w+)\/revision\/(\w+)/
  const matchCommit = result.match(commit)
  if (matchCommit && matchCommit.length > 2) {
    result = result.replace(commit, `<a ${style} href="${matchCommit[0]}">${matchCommit[3]}</a>`)
  }

  return result
}
