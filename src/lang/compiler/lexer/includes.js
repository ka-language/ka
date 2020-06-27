var readfile = require('../imports/read');

global.included = [global.DIRNAME.concat(global.NAME)];

module.exports = (dir, lex) => {

  //loop through lex
  for (let i = 0; i < lex.length; i++)
    if (lex[i].Name == 'include') {
      var include_name = lex[i + 2].Name;

      include_name = include_name.substr(1).slice(0, -1);

      var file
      , sendDir = dir + include_name;

      if (include_name.startsWith('?~')) include_name = './stdlib' + (include_name.startsWith('?~/') ? '' : '/') + include_name.substr(2);

      //if it is a absolute directory, then do not read from the current "dir"
      if (/^[a-zA-Z]\:/.test(include_name)) {
        file = readfile(include_name);
        sendDir = include_name;
      } else file = readfile(dir + include_name);

      if (file.startsWith('Error')) {
        global.errors.push({
          Error: file,
          Dir: dir + global.NAME
        });
        return [];
      }

      for (let o of JSON.parse(file)) {

        if (global.included.includes(o.FileName)) continue;

        global.included.push(o.FileName);

        var lexxed = lexer(o.Content, o.FileName);

        let _lex = lex.slice(0, i)
        , lex_ = lex.slice(i + 3);

        lex = [..._lex, ...lexxed, ...lex_];
      }

    }

  return lex;
};