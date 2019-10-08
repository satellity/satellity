import style from './index.scss';
import React, {Component} from 'react';
import Avatar, {Piece} from 'avataaars';

class Index extends Component {
  constructor(props) {
    super(props);

    this.state = {
      avatar: 'Circle',
      top: 'ShortHairDreads01',
      accessory: 'Prescription02',
      hairColor: 'BrownDark',
      facialHair: 'BeardLight',
      clothe: 'GraphicShirt',
      clotheColor: 'PastelBlue',
      graphic: 'Resist',
      eye: 'Happy',
      eyebrow: 'Default',
      mouth: 'Smile',
      skin: 'Light',
      hair_shapes: ['NoHair','Eyepatch','Hat','Hijab','Turban','WinterHat1','WinterHat2','WinterHat3','WinterHat4','LongHairBigHair','LongHairBob','LongHairBun','LongHairCurly','LongHairCurvy','LongHairDreads','LongHairFrida','LongHairFro','LongHairFroBand','LongHairNotTooLong','LongHairShavedSides','LongHairMiaWallace','LongHairStraight','LongHairStraight2','LongHairStraightStrand','ShortHairDreads01','ShortHairDreads02','ShortHairFrizzle','ShortHairShaggyMullet','ShortHairShortCurly','ShortHairShortFlat','ShortHairShortRound','ShortHairShortWaved','ShortHairSides','ShortHairTheCaesar','ShortHairTheCaesarSidePart'],
      accessories: ['Blank','Kurt','Prescription01','Prescription02','Round','Sunglasses','Wayfarers'],
      facial_hairs: ['Blank','MoustacheFancy','MoustacheMagnum','BeardLight','BeardMedium', 'BeardMajestic'],
      clothes: ['BlazerShirt','BlazerSweater','CollarSweater','GraphicShirt','Hoodie','Overall','ShirtCrewNeck','ShirtScoopNeck','ShirtVNeck'],
      eyes: ['Close','Cry','Default','Dizzy','EyeRoll','Happy','Hearts','Side','Squint','Surprised','Wink','WinkWacky'],
      eyebrows: ['Angry','AngryNatural','Default','DefaultNatural','FlatNatural','RaisedExcited','RaisedExcitedNatural','SadConcerned','SadConcernedNatural','UnibrowNatural','UpDown','UpDownNatural'],
      mouths: ['Concerned','Default','Disbelief','Eating','Grimace','Sad','ScreamOpen','Serious','Smile','Tongue','Twinkle','Vomit'],
      skins: ['Tanned','Yellow','Pale','Light','Brown','DarkBrown','Black'],
      graphics: ['Blank','Skull','SkullOutline','Bat','Cumbia','Deer','Diamond','Hola','Selena','Pizza','Resist','Bear']
    };

    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e, k, v) {
    e.preventDefault();
    this.setState({
      [k]: v
    });
  }

  render () {
    const state = this.state;

    const actions = ['hair','accessory','beard','clothe','eye','eyebrow','mouth', 'skin'].map((o) => {
      return (
        <span key={o} className={style.action}>{o}</span>
      )
    })

    const hairShapes = state.hair_shapes.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'top', o)}>
          <Piece pieceType='top' pieceSize='100' topType={o} hairColor='Blank'/>
        </div>
      )
    });

    const accessories = state.accessories.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'accessory', o)}>
          <Piece pieceType='accessories' pieceSize='100' accessoriesType={o}/>
        </div>
      )
    });

    const facialHairs = state.facial_hairs.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'facialHair', o)}>
          <Piece pieceType='facialHair' pieceSize='100' facialHairType={o}/>
        </div>
      )
    });

    const clothes = state.clothes.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'clothe', o)}>
          <Piece pieceType='clothe' pieceSize='100' clotheType={o}/>
        </div>
      )
    });

    const graphics = state.graphics.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'graphic', o)}>
          <Piece pieceType="graphics" pieceSize="100" graphicType={o} />
        </div>
      )
    });

    const eyes = state.eyes.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'eye', o)}>
          <Piece pieceType='eyes' pieceSize='100' eyeType={o}/>
        </div>
      )
    });

    const eyebrows = state.eyebrows.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'eyebrow', o)}>
          <Piece pieceType='eyebrows' pieceSize='100' eyebrowType={o}/>
        </div>
      )
    });

    const mouths = state.mouths.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'mouth', o)}>
          <Piece pieceType='mouth' pieceSize='100' mouthType={o}/>
        </div>
      )
    });

    const skins = state.skins.map((o) => {
      return (
        <div key={o} className={style.item} onClick={(e) => this.handleClick(e, 'skin', o)}>
          <Piece pieceType='skin' pieceSize='100' skinColor={o}/>
        </div>
      )
    });

    return (
      <div className={style.container}>
        <div className={style.canvas}>
          <div className={style.avatar}>
            <Avatar
              style={{width: '24rem', height: '24rem'}}
              avatarStyle={state.avatar}
              topType={state.top}
              accessoriesType={state.accessory}
              hairColor={state.hairColor}
              facialHairType={state.facialHair}
              clotheType={state.clothe}
              clotheColor={state.clotheColor}
              graphicType={state.graphic}
              eyeType={state.eye}
              eyebrowType={state.eyebrow}
              mouthType={state.mouth}
              skinColor={state.skin} />
          </div>
          <div>
            <div className={style.actions}>
                {actions}
            </div>
            <div>
                {hairShapes}
            </div>
            <div>
                {accessories}
            </div>
            <div>
                {facialHairs}
            </div>
            <div>
                {clothes}
            </div>
            <div>
                {graphics}
            </div>
            <div>
                {eyes}
            </div>
            <div>
                {eyebrows}
            </div>
            <div>
                {mouths}
            </div>
            <div>
                {skins}
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default Index;
