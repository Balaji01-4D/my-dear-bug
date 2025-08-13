import { useMemo, useState } from 'react'
import { Navbar } from '@/components/Navbar'
import { createConfession } from '@/api'
import type { ConfessionRequest } from '@/types'
import { navigate } from '../navigation'

export function SubmitPage() {
	const [form, setForm] = useState<ConfessionRequest>({
		title: '',
		description: '',
		snippet: '',
		language: '',
		tags: [],
		isFlagged: false,
	})
	const [tagsInput, setTagsInput] = useState('')
	const [submitting, setSubmitting] = useState(false)
	const [error, setError] = useState<string | null>(null)
	const [successId, setSuccessId] = useState<number | null>(null)

	const canSubmit = useMemo(() => {
		return (
			form.title.trim().length >= 5 &&
			form.description.trim().length >= 10 &&
			form.language.trim().length > 0 &&
			!submitting
		)
	}, [form, submitting])

	const onSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError(null)
		setSubmitting(true)
		try {
			const payload: ConfessionRequest = {
				...form,
				tags: tagsInput
					.split(',')
					.map(t => t.trim())
					.filter(Boolean),
			}
			const res = await createConfession(payload)
			setSuccessId(res.id)
		} catch (err: any) {
			setError(err?.message || 'Failed to submit')
		} finally {
			setSubmitting(false)
		}
	}

	return (
		<div className="min-h-screen bg-neutral-100 text-black">
			<Navbar />

			<section className="border-b border-neutral-200 bg-white">
				<div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-10">
					<h1 className="text-4xl font-extrabold tracking-tight">Submit your confession</h1>
					<p className="text-neutral-600 mt-2">Share a bug you battled, how you found it, and a code snippet if helpful.</p>
				</div>
			</section>

			<div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8 py-10">
				<div className="bg-white border border-neutral-200 rounded-2xl shadow-soft p-6 sm:p-8">
					{successId ? (
						<div className="text-center">
							<h2 className="text-2xl font-semibold">Thanks for sharing!</h2>
							<p className="text-neutral-600 mt-2">Your confession has been submitted.</p>
							<div className="mt-6 flex items-center justify-center gap-3">
								<button
									onClick={() => navigate('/')}
									className="px-4 py-2 rounded-lg bg-black text-white hover:bg-black/90"
								>
									Go to Home
								</button>
							</div>
						</div>
					) : (
						<form onSubmit={onSubmit} className="space-y-5">
							{error && <div className="rounded-lg border border-red-200 bg-red-50 text-red-700 px-3 py-2">{error}</div>}

							<div>
								<label className="block text-sm font-medium text-neutral-800">Title</label>
								<input
									type="text"
									value={form.title}
									  onChange={e => setForm((f: ConfessionRequest) => ({ ...f, title: e.target.value }))}
									placeholder="A short, descriptive title"
									className="mt-1 w-full h-11 rounded-xl bg-white text-neutral-900 px-3 border border-neutral-300 focus:outline-none focus:border-neutral-400"
									required
									minLength={5}
								/>
								<p className="mt-1 text-xs text-neutral-500">Minimum 5 characters.</p>
							</div>

							<div>
								<label className="block text-sm font-medium text-neutral-800">Language</label>
								<input
									type="text"
									value={form.language}
									  onChange={e => setForm((f: ConfessionRequest) => ({ ...f, language: e.target.value }))}
									placeholder="e.g., go, javascript, python"
									className="mt-1 w-full h-11 rounded-xl bg-white text-neutral-900 px-3 border border-neutral-300 focus:outline-none focus:border-neutral-400"
									required
								/>
							</div>

							<div>
								<label className="block text-sm font-medium text-neutral-800">Description</label>
								<textarea
									value={form.description}
									  onChange={e => setForm((f: ConfessionRequest) => ({ ...f, description: e.target.value }))}
									placeholder="Tell the story: what went wrong, how you debugged, and the fix."
									className="mt-1 w-full rounded-xl bg-white text-neutral-900 px-3 py-2 border border-neutral-300 focus:outline-none focus:border-neutral-400 min-h-[120px]"
									required
									minLength={10}
								/>
								<p className="mt-1 text-xs text-neutral-500">Minimum 10 characters.</p>
							</div>

							<div>
								<label className="block text-sm font-medium text-neutral-800">Code snippet (optional)</label>
								<textarea
									value={form.snippet}
									  onChange={e => setForm((f: ConfessionRequest) => ({ ...f, snippet: e.target.value }))}
									placeholder="Paste any relevant code snippet"
									className="mt-1 w-full rounded-xl bg-white text-neutral-900 px-3 py-2 border border-neutral-300 focus:outline-none focus:border-neutral-400 font-mono text-sm min-h-[140px]"
								/>
							</div>

							<div>
								<label className="block text-sm font-medium text-neutral-800">Tags (optional)</label>
								<input
									type="text"
									value={tagsInput}
									onChange={e => setTagsInput(e.target.value)}
									placeholder="Comma-separated, e.g., concurrency, grpc"
									className="mt-1 w-full h-11 rounded-xl bg-white text-neutral-900 px-3 border border-neutral-300 focus:outline-none focus:border-neutral-400"
								/>
								<p className="mt-1 text-xs text-neutral-500">Add 0-5 tags to help others find this confession.</p>
							</div>

							<div className="pt-2">
								<button
									type="submit"
									disabled={!canSubmit}
									className="px-4 py-2 rounded-lg bg-black text-white font-medium disabled:opacity-50"
								>
									{submitting ? 'Submitting…' : 'Submit Confession'}
								</button>
								<button
									type="button"
									onClick={() => navigate('/')}
									className="ml-3 px-4 py-2 rounded-lg border border-neutral-300"
								>
									Cancel
								</button>
							</div>
						</form>
					)}
				</div>
			</div>
		</div>
	)
}
