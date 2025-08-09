import styled from "styled-components";

export const CollectionsGrid = styled.div`
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 12px;
    padding: 16px;
`;

export const CollectionCard = styled.div`
    display: flex;
    align-items: center;
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    overflow: hidden;
    cursor: pointer;
    transition: transform 0.2s ease;
    height: 64px;

    &:hover {
        transform: scale(1.02);
    }
`;

export const CollectionCover = styled.img`
    width: 64px;
    height: 64px;
    object-fit: cover;
    flex-shrink: 0;
`;

export const CollectionContent = styled.div`
    padding: 0 12px;
    overflow: hidden;
`;

export const CollectionTitle = styled.div`
    font-size: 16px;
    font-weight: 600;
    color: white;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
`;
