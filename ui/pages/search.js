import Head from "next/head";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import useDebounce from "../lib/useDebounce";
import TopNavigation from "../components/top-navigation";
import CodeSearchResult from "../components/code-search-result";
import useSearchCode from "../lib/useSearchCode";

export default function Home() {
  const router = useRouter();
  const { q = "" } = router.query;
  const [notFound, setNotFound] = useState(false);
  const [val, setVal] = useState("");
  const [hits, setHits] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [filter, setFilter] = useState({
    repo: [],
    lang: [],
    path: [],
  });

  const fetchCodeSearch = useSearchCode({ setHits, setNotFound, setIsLoading, setFilter });

  const [, cancel] = useDebounce(
    () => {
      if (val == "") {
        return;
      }
      fetchCodeSearch(val, filter);
    },
    500,
    [val, filter]
  );

  const updateFilter = (filterName, index) => {
    const filter = hits.facets[filterName].buckets[index].val;
    setFilter((prev) => {
      const newFilter = prev;
      if (prev[filterName]?.includes(filter)) {
        const newVal = prev[filterName].filter((item) => item != filter);
        newFilter[filterName] = newVal;
        fetchCodeSearch(val, newFilter);
        return newFilter;
      }

      newFilter[filterName]?.push(filter);
      fetchCodeSearch(val, newFilter);
      return newFilter;
    });
  };

  useEffect(() => {
    if (q !== "" && !hits) {
      setVal(q);
      fetchCodeSearch(q, filter);
    }
  }, [q]);

  useEffect(() => {
    const keyDownHandler = (e) => {
      if (e.code == "Enter") {
        fetchCodeSearch(val, filter);
      }
    };

    document.addEventListener("keydown", keyDownHandler);
    return () => {
      document.removeEventListener("keydown", keyDownHandler);
    };
  }, []);

  return (
    <div className="bg-gray-50 h-screen">
      <Head>
        <meta name="viewport" content="width=device-width,initial-scale=1" />
        <title>{`${q} - heline`}</title>
        <link rel="icon" type="image/png" href="/favicon.png" />
        {/* <script
          defer
          data-domain="heline.dev"
          src="https://plausible.io/js/plausible.js"
        ></script> */}
      </Head>

      <div className="h-full overflow-y-auto scrollbar-hide">
        <nav className="bg-white shadow-sm sticky top-0">
          <div className="py-3 w-full max-w-7xl mx-auto">
            <TopNavigation setVal={setVal} q={q} />
          </div>
        </nav>
        
        <CodeSearchResult
          hits={hits}
          filter={filter}
          updateFilter={updateFilter}
          isLoading={isLoading}
        />

        {notFound && !isLoading && (
          <div className="grid place-items-center pt-32 space-y-8">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-20 w-20 text-blue-500"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
              />
            </svg>
            <div className="text-center text-lg text-gray-600">
              Can not find matching query <strong>"{q}"</strong>.
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
