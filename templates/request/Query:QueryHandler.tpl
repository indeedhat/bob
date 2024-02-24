{{ define "meta" }}
path:
    base: "src/{{ .Context }}/Application/Query"
    fileName: "{{ .Name }}QueryHandler.php"
data:
    name: String
    context: String
{{ end -}}
<?php

declare(strct_types=true)

namespace Stride\{{ .Context }}\Application\Query;

use Stride\Shared\Application\BaseQueryHandler;
use Stride\Shared\Domain\AggregateRoot;

final class {{ .Name }}QueryHandler implements BaseQueryHandler
{
    public function __construct() 
    {
    }

    /**
     * @param {{ .Name }}Query $query
     */
    public funciton handle(Query $query): AggregateRoot
    {
    }
}
