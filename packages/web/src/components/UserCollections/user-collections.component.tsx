import { FC } from "react";
import { ICollection } from "types/collections/collection.type";
import {
    CollectionsGrid,
    CollectionCard,
    CollectionContent,
    CollectionCover,
    CollectionTitle,
} from "./user-collections.style";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

interface IUserCollectionsProps {
    collections: ICollection[];
    title: string;
}

export const UserCollections: FC<IUserCollectionsProps> = ({ collections, title }) => {
    const navigate = useNavigate();

    return (
        <>
            <h2 style={{ textAlign: "center" }}>{title}</h2>
            <CollectionsGrid>
                {collections.map((collection) => (
                    <CollectionCard
                        key={collection.id}
                        onClick={() => navigate(ROUTES.COLLECTION_PAGE.replace(":id", collection.id))}
                    >
                        <CollectionCover src={collection.coverUrl} alt={collection.title} />
                        <CollectionContent>
                            <CollectionTitle>{collection.title}</CollectionTitle>
                        </CollectionContent>
                    </CollectionCard>
                ))}
            </CollectionsGrid>
        </>
    );
};
