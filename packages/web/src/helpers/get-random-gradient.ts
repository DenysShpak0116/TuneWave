export const getRandomGradient = () => {
    const colors = [
        ['#ff7e5f', '#feb47b'],
        ['#6a11cb', '#2575fc'],
        ['#43cea2', '#185a9d'],
        ['#ff4e50', '#f9d423'],
        ['#eecda3', '#ef629f'],
        ['#2193b0', '#6dd5ed'],
        ['#cc2b5e', '#753a88'],
        ['#ee9ca7', '#ffdde1'],
    ];
    const pair = colors[Math.floor(Math.random() * colors.length)];
    return `linear-gradient(135deg, ${pair[0]}, ${pair[1]})`;
};
  