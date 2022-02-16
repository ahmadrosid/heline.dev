import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import useDebounce from '../lib/useDebounce'
import axios from 'axios'
import renderArray from '../lib/render-array'

export default function Home() {
    const router = useRouter()
    const { q = "" } = router.query
    const [notFound, setNotFound] = useState(false)
    const [val, setVal] = useState("")
    const [hits, setHits] = useState(null)
    const [filter, setFilter] = useState({
        repo: [],
        lang: [],
        path: []
    })

    const fetchData = (val, filter) => {
        if (val == "") {
            return
        }
        console.log(filter)

        if (filter == null) {

            if (router.query["filter[repo]"]) {
                const newFilter = Object.assign(filter || {}, {
                    repo: router.query["filter[repo]"].split(",")
                });
                filter = newFilter;
            }

            if (router.query["filter[lang]"]) {
                const newFilter = Object.assign(filter || {}, {
                    lang: router.query["filter[lang]"].split(",")
                });
                filter = newFilter;
            }

            if (router.query["filter[path]"]) {
                const newFilter = Object.assign(filter || {}, {
                    path: router.query["filter[path]"].split(",")
                });
                filter = newFilter;
            }

            if (filter != null) {
                console.log("Updating filter")
                setFilter(() => filter);
            }

        }

        let queryParam = val
        if (filter?.repo?.length > 0) {
            queryParam += `&filter[repo]=${filter.repo.join(",")}`
        }
        
        if (filter?.lang?.length > 0) {
            queryParam += `&filter[lang]=${filter.lang.join(",")}`
        }
        
        if (filter?.path?.length > 0) {
            queryParam += `&filter[path]=${filter.path.join(",")}`
        }
        
        axios.get(`https://heline.dev/api/search?q=${queryParam}`)
        // axios.get(`/api/search?q=${queryParam}`)
            .then(res => {
                if (res.data.hits.hits === null) {
                    setNotFound(true)
                    setHits(null)
                } else {
                    setNotFound(false)
                    setHits(res.data.hits)
                }
            })
            .catch(e => {
                console.error(e)
                setNotFound(true)
                setHits(null)
            })

        if (typeof window !== 'undefined') {
            router.push(`search?q=${queryParam}`)
        }
    }

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

    const updateFilterRepo = (index) => {
        if (hits.facets.repo.buckets.length == 0) {
            return
        }
        updateFilter('repo', index)
    }

    const updateFilterPath = (index) => {
        if (hits.facets.path.buckets.length == 0) {
            return
        }
        updateFilter('path', index)
    }

    const updateFilterLang = (index) => {
        if (hits.facets.lang.buckets.length == 0) {
            return
        }
        updateFilter('lang', index)
    }

    const getPath = (val) => {
        let path = String(val)
        if (path.includes("/")) {
            let paths = path.split("/")
            return paths.slice(Math.max(paths.length - 2, 1)).join("/")
        }

        return val
    }

    useEffect(() => {
        if (q !== "" && !hits) {
            setVal(q)
            fetchData(q, null)
        }
    }, [q])

    return (
        <div className="bg-[#f7f6f3] min-h-screen h-full">
            <Head>
                <meta name='viewport' content='width=device-width,initial-scale=1' />
                <title>{`${q} - heline`}</title>
                <link rel='icon' type='image/png' href='/favicon.png' />
                <script defer data-domain="heline.dev" src="https://plausible.io/js/plausible.js"></script>
            </Head>
            <nav className="border-b bg-white border-green-300">
                <div className="py-6 w-full max-w-7xl mx-auto flex items-center gap-4">
                    <div className="w-full max-w-[25%]">
                        <div className="text-green-500 flex items-center gap-x-2 px-4">
                            <img src="/favicon.png" className="w-8" />
                            <a href="/" className="tracking-wide font-medium uppercase text-3xl pb-1">heline</a>
                        </div>
                    </div>
                    <div className="flex px-4 rounded-lg border border-green-400 bg-white items-center justify-between w-full mr-8 mx-4">
                        <svg viewBox="64 64 896 896" focusable="false" className="text-gray-400" data-icon="search" width="1em" height="1em" fill="currentColor" aria-hidden="true">
                            <path d="M909.6 854.5L649.9 594.8C690.2 542.7 712 479 712 412c0-80.2-31.3-155.4-87.9-212.1-56.6-56.7-132-87.9-212.1-87.9s-155.5 31.3-212.1 87.9C143.2 256.5 112 331.8 112 412c0 80.1 31.3 155.5 87.9 212.1C256.5 680.8 331.8 712 412 712c67 0 130.6-21.8 182.7-62l259.7 259.6a8.2 8.2 0 0 0 11.6 0l43.6-43.5a8.2 8.2 0 0 0 0-11.6zM570.4 570.4C528 612.7 471.8 636 412 636s-116-23.3-158.4-65.6C211.3 528 188 471.8 188 412s23.3-116.1 65.6-158.4C296 211.3 352.2 188 412 188s116.1 23.2 158.4 65.6S636 352.2 636 412s-23.3 116.1-65.6 158.4z"></path>
                        </svg>
                        <input
                            onChange={(el) => {
                                setVal(encodeURIComponent(el.target.value))
                            }}
                            autoFocus={true}
                            spellCheck={false}
                            defaultValue={q}
                            type="text"
                            placeholder="Search"
                            className="text-md bg-white rounded-full w-full focus:outline-none py-3 px-2 text-gray-500"
                        />
                    </div>
                </div>
            </nav>
            {notFound && (
                <div className="grid place-items-center pt-32 space-y-8">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-20 w-20 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    <div className="text-center text-lg text-gray-600">
                        Can not find matching query <strong>"{q}"</strong>.
                    </div>
                </div>
            )}

            {hits && (
                <div className="w-full max-w-7xl mx-auto flex">
                    <div className="w-full min-w-[250px] max-w-[25%] py-8 space-y-4 pl-4">
                        <div className="space-y-2">
                            <h3 className="text-gray-800 uppercase">Repository</h3>
                            {/* <input className="bg-white rounded border border-gray-200 w-full px-4 py-3 text-sm" type="text" placeholder="Filter repositories" /> */}
                            <div className="py-2 space-y-1">
                            {renderArray(
                                hits.facets?.repo?.buckets.map((item, index) => (
                                    <div className="flex justify-between items-center text-gray-600 pr-1">
                                        <div className="flex gap-2 items-center">
                                            <input onChange={() => updateFilterRepo(index)} id={item.val} className="p-2" type="checkbox" checked={filter.repo?.includes(item.val)}/>
                                            <label htmlFor={item.val}>{item.val}</label>
                                        </div>
                                        <div className="text-sm">{item.count}</div>
                                    </div>
                                ))
                            )}
                            </div>
                        </div>
                        <div className="space-y-2 pt-2">
                            <h3 className="text-gray-800 uppercase">Path</h3>
                            {/* <input className="bg-white rounded border border-gray-200 w-full px-4 py-3 text-sm" type="text" placeholder="Filter paths" /> */}
                            <div className="py-2 space-y-1">
                            {renderArray(
                                hits.facets?.path?.buckets.map((item, index) => (
                                    <div className="flex justify-between items-center text-gray-600 pr-1">
                                        <div className="flex gap-2 items-center">
                                            <input onChange={() => updateFilterPath(index)} className="p-2" type="checkbox" checked={filter.path?.includes(item.val)}/>
                                            <label className="truncate">{getPath(item.val)}</label>
                                        </div>
                                        <div className="text-sm">{item.count}</div>
                                    </div>
                                ))
                            )}
                            </div>
                        </div>
                        <div className="space-y-2 pt-2">
                            <h3 className="text-gray-800 uppercase">Language</h3>
                            <div className="py-2 space-y-1">
                            {renderArray(
                                hits.facets?.lang?.buckets.map((item, index) => (
                                    <div className="flex justify-between items-center text-gray-600 pr-1">
                                        <div className="flex gap-2 items-center">
                                            <input onChange={() => updateFilterLang(index)} className="p-2" type="checkbox" checked={filter.lang?.includes(item.val)} />
                                            <label>{item.val}</label>
                                        </div>
                                        <div className="text-sm">{item.count}</div>
                                    </div>
                                ))
                            )}
                            </div>
                        </div>
                    </div>
                    <div className="w-full max-w-[75%] p-8 pr-6">
                        <div className="pb-2">
                            <p className="text-gray-700">Total: {hits.total}</p>
                        </div>
                        {renderArray(
                            hits.hits?.map(item => {
                                if (item.content.snippet === null) {
                                    return;
                                }

                                return (
                                    <div className="py-2">
                                        <div>
                                            <a target="_blank" href={`https://github.com/${item.repo.raw}`} className="flex gap-2 items-center">
                                                <img className="repo-avatar rounded-full border" src={`https://avatars.githubusercontent.com/u/${item.owner_id.raw}?s=60&amp;v=4`} alt="" width="22" height="22" />
                                                <span className="text-gray-700 font-light">{item.repo.raw}</span>
                                            </a>

                                            <a target="_blank" href={`https://github.com/${item.repo.raw}/blob/${item.branch.raw}/${item.file_id.raw.split("/").slice(4, 100).join("/")}`} className="flex gap-1 items-center">
                                                <span className="text-green-500 pl-8">{item.file_id.raw.split("/").slice(2, 100).join("/")}</span>
                                            </a>
                                        </div>
                                        <div className="border border-gray-200 rounded-md bg-white p-2 my-2">
                                            {renderArray(
                                                item.content.snippet?.map((content, parentIndex) => {
                                                    if (content.length === 0) return;

                                                    let contents = []
                                                    let index = 0;
                                                    let chunk = "";
                                                    content.split("\n").forEach((item) => {
                                                        index++;
                                                        chunk = chunk + item;
                                                        if (index === 4) {
                                                            if (chunk.includes("<mark>")) {
                                                                contents.push(chunk);
                                                            }
                                                            chunk = "";
                                                            index = 0;
                                                        }
                                                    })

                                                    if (chunk !== "" && chunk.includes("<mark>")) {
                                                        contents.push(chunk)
                                                    }

                                                    if (contents.length == 0)return;

                                                    // Take only 3 element to render!
                                                    contents = contents.slice(0, 3);

                                                    return (
                                                        <>
                                                            {renderArray(
                                                                contents.map((source, i) => {
                                                                    return (
                                                                        <>
                                                                            <table className="highlight-table">
                                                                                <tbody dangerouslySetInnerHTML={{ __html: source }}></tbody>
                                                                            </table>
                                                                            { (i < contents.length && parentIndex < (item.content.snippet.length - 1)) && <div className="bg-[#f1f0ec] h-8 my-[8px] -mx-2"></div>}
                                                                        </>
                                                                    )
                                                                })
                                                            )}
                                                            {/* { i !== (item.content.snippet.length - 1) && <div className="bg-[#f1f0ec] h-6 my-[8px] -mx-2"></div>} */}
                                                        </>
                                                    )
                                                })
                                            )}

                                        </div>
                                    </div>
                                )
                            })
                        )}
                    </div>
                </div>
            )}
        </div>
    )
}
