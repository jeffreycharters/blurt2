import type { Blurt } from "../../app"
import type { PageServerLoad } from "./$types"

export const load = (async () => {
	return { blurts: [] as Blurt[] }
}) satisfies PageServerLoad
