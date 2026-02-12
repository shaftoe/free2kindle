export async function load({
  parent,
}: {
  parent: () => Promise<{ apiUrl: string }>;
}): Promise<{ apiUrl: string }> {
  const { apiUrl } = await parent();
  return { apiUrl };
}
