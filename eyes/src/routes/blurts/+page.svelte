<script lang="ts">
	import { onDestroy, onMount } from "svelte"
	import { page } from "$app/stores"

	import { Blurt, Loader, Processing } from "$lib/c"

	import { crossfade, fade } from "svelte/transition"
	import { flip } from "svelte/animate"
	import type { PageData } from "./$types"
	const [send, receive] = crossfade({ duration: 200 })

	export let data: PageData
	export let blurts = data.blurts ?? []

	let blurt = ""
	let username = ""
	let initiating = true // loading initial Blurts
	let loading = false // loading new set of Blurts via infinite scroll

	let processing = false // processing adding new Blurt

	onMount(async () => {
		// connect to websocket
	})

	const path = $page.url.origin

	const typeHandler = () => {
		const remaining = 14 - blurt.length
		const hue = (360 - 190) / remaining + 190
		const saturation = (100 - 60) / remaining + 60
		const countdownDiv = document.getElementById("countdown-box")
		if (!countdownDiv) return

		countdownDiv.style.color = `hsl(${hue}, ${saturation}%, 50%)`
		countdownDiv.innerHTML = String(remaining)
	}

	const submitHandler = async () => {
		const countdownDiv = document.getElementById("countdown-box")
		if (!countdownDiv) return

		countdownDiv.style.color = `hsl(190, 60%, 50%)`
		countdownDiv.textContent = ""
	}

	const touchHandler = (e: TouchEvent) => {
		let holderDiv = (e.target as HTMLDivElement).closest("div.blurt-holder") as HTMLDivElement
		if (!holderDiv) return

		holderDiv.style.transform = "scale(0.98)"
	}

	const touchEndHandler = (e: TouchEvent) => {
		let holderDiv = (e.target as HTMLDivElement).closest("div.blurt-holder") as HTMLDivElement
		if (!holderDiv) return

		holderDiv.style.transform = "scale(1)"
		holderDiv.style.boxShadow = ""
	}

	// Function is called to load next set of Blurts by intersection observer.
	// const loadMoreBlurts = async (entries) => {
	// 	entries.forEach(async (e) => {
	// 		if (e.isIntersecting) {
	// 			const loadTimer = setTimeout(() => (loading = true), 500)
	// 			const path = $page.url.origin
	// 			const count = 25
	// 			const url = `${path}/blurts.json?username=${encodeURIComponent(
	// 				username
	// 			)}&take=${count}&cursor=${displayBlurts[displayBlurts.length - 1].uid}`
	// 			const res = await fetch(url)
	// 			if (res.status === 404) {
	// 				document.getElementById("after-blurt").remove()
	// 				return
	// 			}
	// 			const rawBlurts = await res.json()
	// 			const dateBlurts = humanizeDates(rawBlurts)
	// 			blurts.set([...displayBlurts, ...dateBlurts])
	// 			clearTimeout(loadTimer)
	// 			loading = false
	// 		}
	// 	})
	// }

	// if (browser) {
	// 	setTimeout(() => {
	// 		const options = {
	// 			root: null,
	// 			rootMargin: "0px 0px 0px 0px",
	// 			threshold: 1.0
	// 		}
	// 		let observer = new IntersectionObserver(loadMoreBlurts, options)
	// 		let target = document.getElementById("after-blurt")
	// 		observer.observe(target)
	// 	}, 1000)
	// }
</script>

{#if initiating}
	<div
		out:fade|global={{ duration: 200 }}
		class="absolute -top-2 bottom-0 left-0 right-0 z-50 bg-white pt-8 text-center"
	>
		<Processing text="Initiating! Beep Boop..." />
	</div>
{/if}

<div class="mx-auto my-2 max-w-md px-2">
	<form class="flex items-baseline justify-between" on:submit|preventDefault={submitHandler}>
		<label for="blurt" class="hidden">blurt</label>
		<div class="flex items-baseline gap-2">
			<input
				class="w-40 rounded-md border border-solid border-teal-900 bg-teal-50 px-2 py-1 font-semibold text-teal-800"
				type="text"
				bind:value={blurt}
				id="blurt"
				placeholder="blurt here"
				autocomplete="off"
				maxlength="14"
				style="will-change: transform;"
				on:keyup={typeHandler}
			/>
			<div class="font-bold" style="color: hsl(190, 60%, 55%)" id="countdown-box">14</div>
		</div>
		<button
			type="submit"
			id="blurt-button"
			class="rounded-md border border-solid border-gray-900 bg-teal-600 px-4 py-1 font-semibold text-white"
			>blurt it</button
		>
	</form>

	{#if processing}
		<Processing text="Blurting!" />
	{/if}

	<main id="blurt-zone">
		{#each blurts as blurt (blurt.id)}
			<div
				on:touchstart={touchHandler}
				on:touchend={touchEndHandler}
				in:receive|global={{ key: blurt.id }}
				out:send|global={{ key: blurt.id }}
				animate:flip={{ duration: 400 }}
				class="blurt-holder"
			>
				<Blurt {blurt} />
			</div>
		{/each}
		<Loader showing={loading} />
	</main>
	<div id="after-blurt" />
</div>

<style>
	.blurt-holder {
		transition: all 0.075s cubic-bezier(0.17, 0.67, 0.88, 4.06);
	}
</style>
