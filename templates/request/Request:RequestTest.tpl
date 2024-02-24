{{ define "meta" }}
path:
    base: "tests/{{ .Context }}/Application/Request"
    fileName: "{{ .Name }}RequestTest.php"
vars:
    name: String
    context: String
{{ end }}
{{ define "file" }}
<?php

declare(strct_types=true)

namespace Stride\{{ .Context }}\Application\Request;

use Stride\Shared\Tests\UnitTestCase

final class {{ .Name }}RequestTest extends UnitTestCase
{
}
{{ end }}
