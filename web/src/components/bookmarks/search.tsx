import { Button } from "@/components/ui/button";
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from "@/components/ui/command";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import { ChevronDown, Search, X } from "lucide-react";
import { useState } from "react";

import type { BookmarkSearch } from "@/components/bookmarks/types";

export interface SearchToken {
	type: "domain" | "category" | "type";
	value: string;
}

const filterOptions = [
	{
		label: "Types",
		options: [
			"Bookmark",
			"PDF",
			"EPUB",
			"RSS",
			"Newsletter",
			"Image",
			"Video",
			"Podcast",
		],
	},
	{
		label: "Domains",
		options: [],
	},
];

interface SearchTokenProps {
	label: SearchToken["type"];
	value: string;
	onRemove: () => void;
}

interface FilterButtonProps {
	label: string;
	options: string[];
	onSelect: (value: string) => void;
	selectedTokens: SearchToken[]; // Add this new prop
}

const SearchToken: React.FC<SearchTokenProps> = ({
	label,
	value,
	onRemove,
}) => (
	<span className="inline-flex items-center gap-1 px-2 py-1 text-sm bg-blue-50 text-blue-700 rounded-md">
		{label}:<span className="font-medium">{value}</span>
		<X
			className="h-3 w-3 cursor-pointer hover:text-blue-900"
			onClick={onRemove}
		/>
	</span>
);

const FilterButton: React.FC<FilterButtonProps> = ({
	label,
	options,
	onSelect,
	selectedTokens,
}) => {
	const [open, setOpen] = useState(false);
	const [search, setSearch] = useState("");

	const filtered = options.filter((option) =>
		option.toLowerCase().includes(search.toLowerCase()),
	);

	const isOptionSelected = (option: string) =>
		selectedTokens.some(
			(token) =>
				token.type.toLowerCase() === label.toLowerCase() &&
				token.value.toLowerCase() === option.toLowerCase(),
		);

	return (
		<Popover open={open} onOpenChange={setOpen}>
			<PopoverTrigger asChild>
				<Button
					variant="ghost"
					className="h-7 px-2 gap-1 text-gray-600 hover:bg-gray-100 text-sm"
				>
					{label}
					<ChevronDown className="h-3 w-3" />
				</Button>
			</PopoverTrigger>
			<PopoverContent className="w-64 p-0" align="start">
				<Command>
					<CommandInput
						placeholder={`Search ${label.toLowerCase()}...`}
						value={search}
						onValueChange={setSearch}
					/>
					<CommandList>
						<CommandEmpty>No results found.</CommandEmpty>
						<CommandGroup heading={label}>
							{filtered.map((option) => {
								const isSelected = isOptionSelected(option);
								return (
									<CommandItem
										key={option}
										onSelect={() => {
											if (!isSelected) {
												onSelect(option);
												setOpen(false);
											}
										}}
										disabled={isSelected}
										className={
											isSelected ? "opacity-50 cursor-not-allowed" : ""
										}
									>
										<span>{option}</span>
										{isSelected && (
											<span className="ml-2 text-xs text-gray-500">
												(selected)
											</span>
										)}
									</CommandItem>
								);
							})}
						</CommandGroup>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	);
};

interface SearchBoxProps {
	search: BookmarkSearch;
	onSearch?: (tokens: SearchToken[], query: string) => void;
}

const SearchBox: React.FC<SearchBoxProps> = ({ search, onSearch }) => {
	const [tokens, setTokens] = useState<SearchToken[]>(
		search.filters.map((token) => {
			const [type, value] = token.split(":");
			return { type: type as SearchToken["type"], value };
		}),
	);
	const [searchInput, setSearchInput] = useState<string>(search.query);

	const removeToken = (index: number): void => {
		setTokens(tokens.filter((_, i) => i !== index));
	};

	const handleAddToken = (type: SearchToken["type"], value: string): void => {
		// Check if token already exists
		const tokenExists = tokens.some(
			(token) =>
				token.type === type &&
				token.value.toLowerCase() === value.toLowerCase(),
		);

		if (!tokenExists) {
			setTokens([...tokens, { type, value }]);
		}
	};

	const handleSearch = () => {
		onSearch?.(tokens, searchInput);
	};

	return (
		<div className="w-full">
			<div className="w-full space-y-2">
				{/* Main Search Bar */}
				<div className="flex flex-col gap-1">
					<div className="flex justify-between items-center gap-2">
						<div className="relative flex items-center gap-2 w-full border rounded-md pl-3 pr-2 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500">
							<Search className="h-4 w-4 text-gray-400 flex-shrink-0" />
							<div className="flex flex-wrap gap-1 flex-1 py-1.5 max-h-20 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-200 scrollbar-track-transparent">
								{tokens.map((token, index) => (
									<SearchToken
										key={`${token.type}:${token.value}`}
										label={token.type}
										value={token.value}
										onRemove={() => removeToken(index)}
									/>
								))}
								<div className="flex-1 min-w-[150px] flex items-center gap-2">
									<input
										className="w-full outline-none text-sm py-0.5"
										placeholder={tokens.length > 0 ? "" : "Search bookmarks..."}
										value={searchInput}
										onChange={(e) => setSearchInput(e.target.value)}
										onKeyDown={(e) => {
											if (e.key === "Enter") {
												handleSearch();
											}
										}}
									/>

									{tokens.length > 0 && (
										<span className="text-xs text-gray-500 px-1.5 py-0.5 bg-gray-100 rounded-full">
											{tokens.length} filter{tokens.length !== 1 ? "s" : ""}
										</span>
									)}
								</div>
							</div>
						</div>
						<Button type="button" onClick={handleSearch}>
							Search
						</Button>
					</div>

					{/* Filters Bar */}
					<div className="flex items-center gap-1 bg-white">
						<span className="text-xs text-gray-500">Filters:</span>
						<div className="flex flex-wrap gap-1">
							{filterOptions.map((filter) => (
								<FilterButton
									key={filter.label}
									label={filter.label}
									options={filter.options}
									onSelect={(value) =>
										handleAddToken(
											filter.label.toLowerCase() as SearchToken["type"],
											value,
										)
									}
									selectedTokens={tokens}
								/>
							))}
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default SearchBox;
