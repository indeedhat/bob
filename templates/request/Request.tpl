{{ define "meta" }}
path:
    base: "src/{{ .Context }}/Application/Request"
    fileName: "{{ .Name }}Request.php"
data:
    name: String
    context: String
{{ end -}}
<?php

declare(strct_types=true)

namespace Stride\{{ .Context }}\Application\Request;

use Stride\Shared\Application\BaseRequest;

final class {{ .Name }}Request extends BaseRequest
{
}
