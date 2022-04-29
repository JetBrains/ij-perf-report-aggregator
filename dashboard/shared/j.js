/////////////////////////////////////////////////////////////////////////////////
// JSONCrush v1.1.6 by Frank Force - https://github.com/KilledByAPixel/JSONCrush
/////////////////////////////////////////////////////////////////////////////////

'use strict';

// create a string of replacement characters
   let characters = [];

   // prefer replacing with characters that will not be escaped by encodeURIComponent
   const unescapedCharacters = `-_.!~*'()`;
   for (let i=127; --i;)
   {
       if
       (
           (i>=48 && i<=57) || // 0-9
           (i>=65 && i<=90) || // A-Z
           (i>=97 && i<=122)|| // a-z
           unescapedCharacters.includes(String.fromCharCode(i))
       )
           characters.push(String.fromCharCode(i));
   }

   // pick from extended set last
   for (let i=32; i<255; ++i)
   {
       let c = String.fromCharCode(i);
       if (c!=='\\' && !characters.includes(c))
           characters.unshift(c);
   }

module.exports  =
{

crush: (string, maxSubstringLength=50)=>
{
    const delimiter = '\u0001'; // used to split parts of crushed string
    const JSCrush=(string, replaceCharacters)=>
    {
        // JSCrush Algorithm (repleace repeated substrings with single characters)
        let replaceCharacterPos = replaceCharacters.length;
        let splitString = '';

        const ByteLength =(string)=>encodeURI(encodeURIComponent(string)).replace(/%../g,'i').length;
        const HasUnmatchedSurrogate =(string)=>
        {
            // check ends of string for unmatched surrogate pairs
            let c1 = string.charCodeAt(0);
            let c2 = string.charCodeAt(string.length-1);
            return (c1 >= 0xDC00 && c1 <= 0xDFFF) || (c2 >= 0xD800 && c2 <= 0xDBFF);
        }

        // count instances of substrings
        let substringCount = new Map();
        for (let substringLength = 2; substringLength < maxSubstringLength; substringLength++)
        for (let i = 0; i < string.length - substringLength; ++i)
        {
            let substring = string.substr(i, substringLength);

            // don't recount if already in list
            if (substringCount.has(substring))
                continue;

            // prevent breaking up unmatched surrogates
            if (HasUnmatchedSurrogate(substring))
                continue;

            // count how many times the substring appears
            let count = 1;
            for (let substringPos = string.indexOf(substring, i+substringLength); substringPos >= 0; ++count)
                substringPos = string.indexOf(substring, substringPos + substringLength);

            // add to list if it appears multiple times
            if (count > 1)
                substringCount.set(substring, count);
        }

        while(true) // loop while string can be crushed more
        {
            // get the next character that is not in the string
            for (;replaceCharacterPos-- && string.includes(replaceCharacters[replaceCharacterPos]);){}
            if (replaceCharacterPos < 0)
                break; // ran out of replacement characters
            let replaceCharacter = replaceCharacters[replaceCharacterPos];

            // find the longest substring to replace
            let bestSubstring;
            let bestLengthDelta = 0;
            let replaceByteLength = ByteLength(replaceCharacter);
            for (let substring of substringCount.keys())
            {
                // calculate change in length of string if it substring was replaced
                let count = substringCount.get(substring);
                let lengthDelta = (count-1)*ByteLength(substring) - (count+1)*replaceByteLength;
                if (!splitString.length)
                    lengthDelta -= ByteLength(delimiter); // include the delimiter length
                if (lengthDelta <= 0)
                  substringCount.delete(substring)
                else if (lengthDelta > bestLengthDelta)
                {
                    bestSubstring = substring
                    bestLengthDelta = lengthDelta;
                }
            }
            if (!bestSubstring)
                break; // string can't be compressed further

            // create new string with the split character
            string = string.split(bestSubstring).join(replaceCharacter) + replaceCharacter + bestSubstring;
            splitString = replaceCharacter + splitString;

            // update substring count list after the replacement
            let newSubstringCount = new Map();
            for (let substring of substringCount.keys())
            {
                // make a new substring with the replacement
                let newSubstring = substring.split(bestSubstring).join(replaceCharacter);

                // count how many times the new substring appears
                let count = 0;
                for (let i = string.indexOf(newSubstring); i >= 0; ++count)
                    i = string.indexOf(newSubstring, i + newSubstring.length);

                // add to list if it appears multiple times
                if (count > 1)
                    newSubstringCount.set(newSubstring, count);

            }
            substringCount = newSubstringCount;
        }

        return {a:string, b:splitString};
    }



    // remove delimiter if it is found in the string
    string = string.replace(new RegExp(delimiter,'g'),'');

    // swap out common json characters

    // crush with JS crush
    const crushed = JSCrush(string, characters);

    // insert delimiter between JSCrush parts
    let crushedString = crushed.a;
    if (crushed.b.length)
        crushedString += delimiter + crushed.b;

    // fix issues with some links not being recognized properly
    crushedString += '_';

    // return crushed string
    return crushedString;
},

uncrush: (string)=>
{
    // remove last character
    string = string.substring(0, string.length - 1);

    // unsplit the string using the delimiter
    const stringParts = string.split('\u0001');

    // JSUncrush algorithm
    let uncrushedString = stringParts[0];
    if (stringParts.length > 1)
    {
        let splitString = stringParts[1];
        for (let character of splitString)
        {
            // split the string using the current splitCharacter
            let splitArray = uncrushedString.split(character);

            // rejoin the string with the last element from the split
            uncrushedString = splitArray.join(splitArray.pop());
        }
    }

    // unswap the json characters in reverse direction
    return uncrushedString;
}

} // JSONCrush
