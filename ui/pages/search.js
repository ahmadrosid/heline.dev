import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import useDebounce from '../lib/useDebounce'
import SubNavigation from "../components/sub-navigation"
import TopNavigation from '../components/top-navigation'
import CodeSearchResult from '../components/code-search-result'
import useSearchCode from '../lib/useSearchCode'

export default function Home() {
    const router = useRouter()
    const { q = "", tbm = "" } = router.query
    const [notFound, setNotFound] = useState(false)
    const [val, setVal] = useState("")
    const [hits, setHits] = useState(null)
    const [filter, setFilter] = useState({
        repo: [],
        lang: [],
        path: []
    })

    const fetchData = useSearchCode({ setHits, setNotFound })

    const [, cancel] = useDebounce(
        () => {
            if (val == '') {
                return;
            }
            fetchData(val, filter)
        },
        500,
        [val]
    )

    const updateFilter = (filterName, index) => {
        const filter = hits.facets[filterName].buckets[index].val
        setFilter(prev => {
            const newFilter = prev
            if (prev[filterName]?.includes(filter)) {
                const newVal = prev[filterName].filter(item => item != filter)
                newFilter[filterName] = newVal
                fetchData(val, newFilter)
                return newFilter
            }

            newFilter[filterName]?.push(filter)
            fetchData(val, newFilter)
            return newFilter
        })
    }

    const updateMatchingSearch = (val) => {
        const { pathname, query } = router
        query.tbm = val;
        router.push({ pathname, query });
        console.log(query)
    }

    useEffect(() => {
        if (q !== "" && !hits) {
            setVal(q)
            fetchData(q, null)
        }
    }, [q])

    return (
        <div className="bg-zinc-50 min-h-screen h-full">
            <Head>
                <meta name='viewport' content='width=device-width,initial-scale=1' />
                <title>{`${q} - heline`}</title>
                <link rel='icon' type='image/png' href='/favicon.png' />
                <script defer data-domain="heline.dev" src="https://plausible.io/js/plausible.js"></script>
            </Head>

            <nav className="bg-white shadow-sm">
                <div className="pt-6 w-full max-w-7xl mx-auto">
                    <TopNavigation setVal={setVal} q={q} />

                    <SubNavigation updateMatchingSearch={updateMatchingSearch} tbm={tbm} />
                </div>
            </nav>

            {notFound && (
                <div className="grid place-items-center pt-32 space-y-8">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-20 w-20 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    <div className="text-center text-lg text-gray-600">
                        Can not find matching query <strong>"{q}"</strong>.
                    </div>
                </div>
            )}

            <CodeSearchResult hits={hits} filter={filter} updateFilter={updateFilter} />
        </div>
    )
}
