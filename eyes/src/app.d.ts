// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export type User = {
	id: Uuid
	username: string
	createdAt: Date
}

export type Blurt = {
	id: Uuid
	user: Uuid
	created: Date
	content: string
	liks: User[]
}

export type Lik = {
	user: Uuid
	blurt: Uuid
	created: Date
}
