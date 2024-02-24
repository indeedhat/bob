{{ define "meta" }}
path:
    base: "src/{{ index . "context" }}/Application/Query"
    fileName: "{{ index . "name" }}Query.php"
vars:
    name: String
    context: String
    params: Map
{{ end }}
{{ define "file" }}
<?php

declare(strct_types=true)

namespace Stride\{{ index . "context" }}\Application\Query;

use Stride\Shared\Application\BaseQuery;
use Stride\Shared\Application\Dto;

final class {{ index . "name" }}Query implements BaseQuery
{
    public function __construct(
        {{- range $value, $key := (index . "params") -}}
            private {{ $value }} ${{ $key }},
        {{ end -}}
    ) {
    }

    public static funciton fromDto(Dto $dto): self
    {
        return new self(
            {{- range (index . "params") -}}
                $dto->{{ . }},
            {{ end -}}
        );
    }

    {{- range $value, $key := (index . "params") -}}
        public function get{{ $key | ucfirst }}(): {{ $value }}
        {
            return $this->{{ $value }};
        }
    {{ end }}
}
{{ end }}
