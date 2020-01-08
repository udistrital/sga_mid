# sga_mid

... ... ...

## Especificación Técnica

Mid del SGA

### Requirements
Go version >= 1.8.

### Preparation:

Para usar el API, usar el comando:

```bash
go get github.com/udistrital/sga_mid
```
### Run

```bash
... ...
```

### Run Tests

```bash
... ...
```
## Example:
```bash
PERSONAS_SERVICE=localhost:8083/v1 CORE_SERVICE=localhost:8102/v1 PROGRAMA_ACADEMICO_SERVICE=localhost:8101/v1 EVENTOS_SERVICE:localhost:8083/v1 OIKOS_SERVICE:10.20.0.254/oikos_api/v1 SGA_MID_HTTP_PORT=8095 bee run -downdoc=true -gendoc=true
```

## Licencia

This file is part of sga_mid.

sga_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

Foobar is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Foobar. If not, see https://www.gnu.org/licenses/.
