import { FC } from "react";
import { Container, Description, Info, Title, Image, StyledRadio } from "./collection-select-item.style";
import { ICollection } from "types/collections/collection.type";

export type IGetCollectionResponse = Omit<ICollection, "collectionSongs" | "user">

interface CollectionSelectItemProps {
    collection: IGetCollectionResponse
    selected: boolean;
    onSelect: () => void;
}

export const CollectionSelectItem: FC<CollectionSelectItemProps> = ({
    collection,
    selected,
    onSelect
}) => {
    return (
        <label>
            <Container>
                <Image src={collection.coverUrl} alt="cover" />
                <Info>
                    <Title>{collection.title}</Title>
                    <Description>{collection.description}</Description>
                </Info>
                <StyledRadio
                    type="radio"
                    checked={selected}
                    onChange={onSelect}
                    name="collection"
                    value={collection.id}
                />
            </Container>
        </label>
    );
};