export const Routes = await fetch(new URL('../routes.json', import.meta.url)).then((res) => res.json())
