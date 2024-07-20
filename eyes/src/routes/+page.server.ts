import type { PageServerLoad } from "./$types"
import { z } from "zod"
import { message, superValidate } from "sveltekit-superforms"
import { zod } from "sveltekit-superforms/adapters"
import { fail } from "@sveltejs/kit"

const schema = z.object({
	username: z
		.string()
		.min(1, { message: "Username please!" })
		.max(14, { message: "Username too long! Max 14 characters" })
})

export const load = (async () => {
	const form = await superValidate(zod(schema))
	return { form }
}) satisfies PageServerLoad

export const actions = {
	default: async ({ request }) => {
		const form = await superValidate(request, zod(schema))
		if (!form.valid) {
			return fail(400, { form })
		}

		try {
			const res = await fetch(`http://localhost:3320/api/v1/test`)
		} catch (err) {
			return message(form, "Login failed! Try later.", { status: 500 })
		}

		return { form }
	}
}
