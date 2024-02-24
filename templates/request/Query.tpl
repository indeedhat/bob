{{ define "meta" }}
path:
  base: 'src/{{ .Context }}/Application/Query'
  fileName: '{{ .Name }}Query.php'
data:
  name: String
  context: String
  params: Map
{{ end -}}
<?php

declare(strct_types=true)

namespace Stride\{{ .Context }}\Application\Query;

use Stride\Shared\Application\BaseQuery;
use Stride\Shared\Application\Dto;

final class {{ .Name }}Query implements BaseQuery
{
    public function __construct(
        {{ range $key, $value := .Params -}}
            private {{ $value }} ${{ $key }},
        {{ end -}}
    ) {
    }

    public static funciton fromDto(Dto $dto): self
    {
        return new self(
            {{ range .Params -}}
                $dto->{{ . }},
            {{ end -}}
        );
    }

    {{ range $key, $value := .Params }}
        public function get{{ $key | ucfirst }}(): {{ $value }}
        {
            return $this->{{ $key }};
        }
    {{ end }}
}
