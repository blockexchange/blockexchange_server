package types_test

import (
	"blockexchange/types"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalMediafile(t *testing.T) {
	str := `{"data":"IyBCbGVuZGVyIHYyLjc3IChzdWIgMCkgT0JKIEZpbGU6ICd0b3JjaF9jZWlsaW5nLmJsZW5kJwojIHd3dy5ibGVuZGVyLm9yZwp2IC0wLjA2MjQ2OSAtMC4wNDczMzEgMC4wNjgxNTIKdiAtMC4wNjI0NjkgLTAuNTU5NTE1IC0wLjE2NDM4OAp2IC0wLjA2MjQ2OSAwLjAwNDM0NCAtMC4wNDU2NjcKdiAtMC4wNjI0NjkgLTAuNTA3ODM5IC0wLjI3ODIwNgp2IDAuMDYyNTMxIC0wLjA0NzMzMSAwLjA2ODE1Mgp2IDAuMDYyNTMxIC0wLjU1OTUxNSAtMC4xNjQzODgKdiAwLjA2MjUzMSAwLjAwNDM0NCAtMC4wNDU2NjcKdiAwLjA2MjUzMSAtMC41MDc4MzkgLTAuMjc4MjA2CnYgMC4zNTM1ODQgMC4wNDAwMDAgMC4zNjM1NTMKdiAwLjM1MzU4NCAtMC4zOTc1MDAgMC4zNjM1NTMKdiAtMC4zNTM1MjIgMC4wNDAwMDAgLTAuMzQzNTUzCnYgLTAuMzUzNTIyIC0wLjM5NzUwMCAtMC4zNDM1NTMKdiAwLjM1MzU4NCAwLjA0MDAwMCAtMC4zNDM1NTMKdiAtMC4zNTM1MjIgMC4wNDAwMDAgMC4zNjM1NTMKdiAwLjM1MzU4NCAtMC4zOTc1MDAgLTAuMzQzNTUzCnYgLTAuMzUzNTIyIC0wLjM5NzUwMCAwLjM2MzU1Mwp2dCAwLjU2MjUgMC41MDAwCnZ0IDAuNTYyNSAwLjYyNTAKdnQgMC40Mzc1IDAuNjI1MAp2dCAwLjQzNzUgMC41MDAwCnZ0IDAuNDM3NSAwLjAwMDAKdnQgMC41NjI1IDAuMDAwMAp2dCAwLjU2MjUgMC4xMjUwCnZ0IDAuNDM3NSAwLjEyNTAKdnQgMC41NjI1IDAuNjI1MAp2dCAwLjQzNzUgMC42MjUwCnZ0IDAuNDM3NSAwLjYyNTAKdnQgMC40Mzc1IDAuMDAwMAp2dCAwLjU2MjUgMC42MjUwCnZ0IDAuNTYyNSAwLjAwMDAKdnQgMS4wMDAwIDAuNTYyNQp2dCAxLjAwMDAgMS4wMDAwCnZ0IDAuMDAwMCAxLjAwMDAKdnQgMC4wMDAwIDAuNTYyNQp2dCAwLjAwMDAgMC41NjI1CnZ0IDEuMDAwMCAwLjU2MjUKdnQgMS4wMDAwIDEuMDAwMAp2dCAwLjAwMDAgMS4wMDAwCnZuIDAuMDAwMCAwLjkxMDUgMC40MTM0CnZuIC0wLjAwMDAgLTAuNDEzNCAwLjkxMDUKdm4gLTEuMDAwMCAwLjAwMDAgMC4wMDAwCnZuIDAuNzA3MSAwLjAwMDAgLTAuNzA3MQp2biAwLjcwNzEgMC4wMDAwIDAuNzA3MQpmIDMvMS8xIDEvMi8xIDUvMy8xIDcvNC8xCmYgOC81LzEgNC82LzEgMi83LzEgNi84LzEKZiAzLzkvMiA0LzYvMiA4LzUvMiA3LzEwLzIKZiAxLzExLzMgMy85LzMgNC82LzMgMi8xMi8zCmYgNS8xMy8yIDEvMTEvMiAyLzEyLzIgNi8xNC8yCmYgNy8xMC8zIDgvNS8zIDYvMTQvMyA1LzEzLzMKZiA5LzE1LzQgMTAvMTYvNCAxMi8xNy80IDExLzE4LzQKZiAxMy8xOS81IDE0LzIwLzUgMTYvMjEvNSAxNS8yMi81Cg","mod_name":"default","name":"torch_ceiling.obj"}`

	mf := &types.Mediafile{}
	err := json.Unmarshal([]byte(str), mf)
	assert.NoError(t, err)
	assert.NotNil(t, mf.Data)
}