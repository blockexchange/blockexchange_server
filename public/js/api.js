
export const browse = path => m.request({
    method: "POST",
    url: "api/browse",
    data: {
			path: path
		}
});
