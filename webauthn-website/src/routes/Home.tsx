import React from 'react';
import chunk from 'lodash.chunk';
import { Client } from '../http';

type Props = {
  client: Client;
};

type Place = {
  id: string;
  name: string;
  imageURL: string;
};

const Home: React.FC<Props> = ({ client }) => {
  const [places, setPlaces] = React.useState<Array<Place[]>>([]);

  React.useEffect(() => {
    let isCancelled = false;
    client.get<{ places: Place[] }>('/list').then(response => {
      if (!isCancelled) {
        setPlaces(chunk(response.places, 2));
      }
    });

    return () => {
      isCancelled = true;
    };
  }, [client]);

  return (
    <>
      {places.map((_, idx) => (
        <div key={idx} className="tile is-ancestor">
          {_.map(place => (
            <div key={place.id} className="tile is-parent">
              <article className="tile is-child box">
                <p className="title">{place.name}</p>
                <figure className="image is-4by3">
                  <img alt={place.name} src={place.imageURL} />
                </figure>
              </article>
            </div>
          ))}
        </div>
      ))}
    </>
  );
};

export default Home;
