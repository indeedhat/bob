# Bob (the builder)

Bob is a tool for generating boilerplate files for your project

## Define templates
Templates should be defined in the templates directory (for now).  
The location of the file does not matter but the format of the name does `[groupName]:[templateName].tpl`. If the file name is defined without the `:` then the group and file names will both be set to the full string.

Templates are handled using the text/template standard library.

Each template must define a meta block:
```tpl
{{ define "meta" }}
# all entries in the path map can contain any data replacements as defined in the data section
path:
  # base defines the directory that the file will be saved in
  base: 'src/{{ .Context }}/Application/Query' 
  # fileName defines the final name of the file to be saved adter the template is executed
  fileName: '{{ .Name }}Query.php'

# the data  map defines the shape of the expected input data
# Note: when using the data from a template each of the keys in this map will be transformed with ucfirst
data:
  # defines the key and expected data type for the input data
  # possible data types include: 
  # - String (a basic string literal)
  # - List (a list of type []string)
  # - Map (a map of type map[string]string)
  name: String
  context: String
  params: Map
{{ end -}}
{{/* Everything after this block contains the actual file template */}}
}
```

## TODO
- [ ] decide how i want bob to pickup templates
    - baked in?
    - current dir?
    - .config dir?
- [ ] pretty up the ui a little
    - [ ] Better output for maker panel
    - [ ] side by side key/value fields for map
- [ ] make reverting partial successes an option
- [ ] newline on enter for text boxes
