{{ define "meta" }}
path:
    base: "tests/{{ .Context }}/Application/Request"
    fileName: "{{ .Name }}RequestTest.php"
data:
    name: String
    context: String
{{ end -}}
<?php

declare(strct_types=true)

namespace Stride\{{ .Context }}\Application\Request;

use Stride\Shared\Tests\UnitTestCase

final class {{ .Name }}RequestTest extends UnitTestCase
{
}
