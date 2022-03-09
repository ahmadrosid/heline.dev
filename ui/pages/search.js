import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import useDebounce from '../lib/useDebounce'
import SubNavigation from "../components/sub-navigation"
import TopNavigation from '../components/top-navigation'
import CodeSearchResult from '../components/code-search-result'
import useSearchCode from '../lib/useSearchCode'
import useSearchDocument from '../lib/useSearchDocument'
import { IoRocketOutline } from "react-icons/io5"
import DocSearchResult from '../components/docs-search-result'

export default function Home() {
    const router = useRouter()
    const { q = "", tbm = "" } = router.query
    const [notFound, setNotFound] = useState(false)
    const [val, setVal] = useState("")
    const [hits, setHits] = useState(null)
    const [hitsDocs, setDocsHits] = useState(null)
    const [isLoading, setIsLoading] = useState(false)
    const [filter, setFilter] = useState({
        repo: [],
        lang: [],
        path: []
    })

    const fetchCodeSearch = useSearchCode({ setHits, setNotFound, setIsLoading })
    const fetchDocumentSearch = useSearchDocument({ setHits: setDocsHits, setNotFound, setIsLoading })

    const [, cancel] = useDebounce(
        () => {
            if (val == '') {
                return;
            }
            if (tbm === "code" || tbm === "") {
                fetchCodeSearch(val, filter)
            }
            if (tbm === "docs") {
                fetchDocumentSearch(val)
            }
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
                fetchCodeSearch(val, newFilter)
                return newFilter
            }

            newFilter[filterName]?.push(filter)
            fetchCodeSearch(val, newFilter)
            return newFilter
        })
    }

    const updateMatchingSearch = (tbm) => {
        const { pathname, query } = router
        query.tbm = tbm;
        router.push({ pathname, query });
        if (hitsDocs === null) {
            fetchDocumentSearch(val)
        }
        if (hits === null) {
            fetchCodeSearch(val)
        }
    }

    useEffect(() => {
        if (q !== "" && !hits) {
            setVal(q)
            if (tbm === "code" || tbm === "") {
                fetchCodeSearch(val, null)
            }
            if (tbm === "docs") {
                fetchDocumentSearch(val)
            }
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

            {(tbm === "code" || tbm === "") && (
                <CodeSearchResult
                    hits={hits}
                    filter={filter}
                    updateFilter={updateFilter}
                    isLoading={isLoading}
                />
            )}

            {(tbm === "docs" ) && (
                <DocSearchResult
                    hits={hitsDocs}
                    isLoading={isLoading}
                />
            )}

            {(notFound && !isLoading) && (
                <div className="grid place-items-center pt-32 space-y-8">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-20 w-20 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    <div className="text-center text-lg text-gray-600">
                        Can not find matching query <strong>"{q}"</strong>.
                    </div>
                </div>
            )}

            {(tbm === "docsx" || tbm === "stf" || tbm === "blog") && (
                <div className="grid place-items-center pt-32 space-y-4">
                    <IoRocketOutline className='text-9xl text-emerald-400' />
                    <div className="text-center text-4xl text-gray-600">
                        <span className='font-medium'>Coming soon!</span>
                    </div>
                </div>
            )}
        </div>
    )
}
