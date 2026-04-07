INSERT INTO keyboard_zones (name, symbols)
VALUES
    ('en_red', '`,1,q,a,z,Tab,CapsLock,Shift'),
    ('en_orange', '2,w,s,x'),
    ('en_yellow', '3,e,d,c'),
    ('en_green', '4,5,r,t,f,g,v,b'),
    ('en_blue', '6,7,y,u,h,j,n,m'),
    ('en_indigo', '8,i,k,Comma'),
    ('en_purple', '9,o,l,Period'),
    ('en_pink', '0,p,Semicolon,Slash,Minus,Equals,Backspace,Backslash,Enter,Shift'),
    ('en_thumb', 'Space')
ON CONFLICT (name) DO UPDATE
SET symbols = EXCLUDED.symbols;

